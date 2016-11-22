package db

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	k8sTokenPath               = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	envVarK8SPort              = "KUBERNETES_PORT_443_TCP_PORT"
	envVarK8SAddress           = "KUBERNETES_PORT_443_TCP_ADDR"
	httpHeaderAccept           = "Accept"
	envVarGophrDBService       = "GOPHR_DB_ADDR"
	httpHeaderBearerPrefix     = "Bearer "
	httpHeaderAuthorization    = "Authorization"
	getK8SEndpointsURLTemplate = "https://%s:%s/api/v1/namespaces/gophr/endpoints/%s"
)

// k8sEndpoints is the structure of a k8s API endpoints request.
type k8sEndpoints struct {
	Subsets []struct {
		Addresses []struct {
			IP string `json:"ip"`
		} `json:"addresses"`
	} `json:"subsets"`
}

// getDBNodes gets the IP addresses of the database nodes in the cluster.
func getDBNodes() ([]string, error) {
	var (
		err            error
		k8sPort        = os.Getenv(envVarK8SPort)
		rawToken       []byte
		k8sAddress     = os.Getenv(envVarK8SAddress)
		gophrDBService = os.Getenv(envVarGophrDBService)
	)

	// First, validate the environment variables.
	if len(k8sPort) < 1 {
		return nil, fmt.Errorf(
			`Couldn't get db nodes: $%s was empty.`,
			envVarK8SPort)
	} else if len(k8sAddress) < 1 {
		return nil, fmt.Errorf(
			`Couldn't get db nodes: $%s was empty.`,
			envVarK8SAddress)
	} else if len(gophrDBService) < 1 {
		return nil, fmt.Errorf(
			`Couldn't get db nodes: $%s was empty.`,
			envVarGophrDBService)
	}

	// Then read the token.
	if rawToken, err = ioutil.ReadFile(k8sTokenPath); err != nil {
		return nil, fmt.Errorf(
			`Couldn't get db nodes: Couldn't parse the request URL: %v.`,
			err)
	}

	// Create the pre-requisites for a new request.
	var (
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client             = &http.Client{Transport: tr}
		payload            k8sEndpoints
		getK8SEndpointsURL = fmt.Sprintf(
			getK8SEndpointsURLTemplate,
			k8sAddress,
			k8sPort,
			gophrDBService)
		getK8SEndpointsRequest  *http.Request
		getK8SEndpointsResponse *http.Response
	)

	// Create the request.
	if getK8SEndpointsRequest, err = http.NewRequest(
		http.MethodGet,
		getK8SEndpointsURL,
		nil,
	); err != nil {
		return nil, fmt.Errorf(
			`Couldn't get db nodes: Failed to create the request: %v.`,
			err)
	}

	// Add the authorization header.
	getK8SEndpointsRequest.Header.Add(httpHeaderAccept, "*/*")
	getK8SEndpointsRequest.Header.Add(
		httpHeaderAuthorization,
		httpHeaderBearerPrefix+string(rawToken[:]))

	// Submit the request.
	if getK8SEndpointsResponse, err = client.Do(
		getK8SEndpointsRequest,
	); err != nil {
		return nil, fmt.Errorf(
			`Couldn't get db nodes: Failed to communicate with K8S: %v.`,
			err)
	} else if getK8SEndpointsResponse.StatusCode >= 400 {
		return nil, fmt.Errorf(
			`Couldn't get db nodes: Failed to communicate with K8S: `+
				`Bumped into status %d.`,
			getK8SEndpointsResponse.StatusCode)
	}

	// Read the body of the request.
	if err = json.
		NewDecoder(getK8SEndpointsResponse.Body).
		Decode(&payload); err != nil {
		return nil, fmt.Errorf(
			`Couldn't get db nodes: Failed to unmarshal the K8S response: %v.`,
			err)
	}

	// Valiate the payload.
	if len(payload.Subsets) < 1 {
		return nil, fmt.Errorf(
			`Couldn't get db nodes: Failed to unmarshal the K8S response: %s.`,
			`there weren't any subsets`)
	} else if len(payload.Subsets[0].Addresses) < 1 {
		return nil, fmt.Errorf(
			`Couldn't get db nodes: Failed to unmarshal the K8S response: %s.`,
			`there weren't any addresses`)
	}

	// Read the IPs from the payload.
	var IPs []string
	for _, address := range payload.Subsets[0].Addresses {
		IPs = append(IPs, address.IP)
	}

	return IPs, nil
}
