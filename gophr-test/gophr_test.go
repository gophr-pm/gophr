package gophr_test

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/lib/db/query"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os/exec"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const (
	baseGophrUrl    = "https://%s:30443/%s%s?go-get=1"
	packageRelation = "gophr.package_archive_records"
)

var (
	expectedShas = [4]string{"3726a1196606290d81728d22412586dc6b2e0327", "9bb4a68d57ff6f623363aa172f0a8297aa289ba7", "c4939d1166b2220bb45338e21506623b4bbdec50", "69483b4bd14f5845b5a1e55bca19e954e827f1d0"}
	repos        = [4]string{"linkedin/burrow", "Shopify/sarama", "uber-go/zap", "stretchr/testify"}
	semVers      = [4]string{"@3726a1", "@v1.10.0", "@c4939d1166b2220bb45338e21506623b4bbdec50", "@69483b4bd14f5845b5a1e55bca19e954e827f1d0/assert"}
)

type RunInfo struct {
	success bool
	message string
}

type TestCase struct {
	gophrUrl   string
	minikubeIp string
	httpClient *http.Client
	runCh      chan RunInfo
	runs       int64
	session    *gocql.Session
	testing    *testing.T
}

func newTestCase(t *testing.T) *TestCase {
	cmd, err := exec.Command("minikube", "ip").CombinedOutput()
	if err != nil {
		panic("Minikube operation failed " + err.Error())
	}
	minikubeIp := strings.Trim(string(cmd), "\n")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	cluster := gocql.NewCluster(minikubeIp + ":30942")
	cluster.ProtoVersion = query.DBProtoVersion
	cluster.Consistency = gocql.One
	session, err := cluster.CreateSession()
	if err != nil {
		panic("Failed to create session " + err.Error())
	}

	err = query.CreateKeyspaceIfNotExists().
		WithReplication("SimpleStrategy", 1).
		WithDurableWrites(true).
		Create(session).
		Exec()
	if err != nil {
		panic("Failed to create keyspace " + err.Error())
	}

	return &TestCase{
		minikubeIp: minikubeIp,
		httpClient: client,
		runCh:      make(chan RunInfo),
		session:    session,
		runs:       0,
		testing:    t,
	}
}

// set up DB, clean stuff, whatever
func (t *TestCase) setup() {
	err := t.session.Query("TRUNCATE " + packageRelation).Exec()
	if err != nil {
		panic("Failure to truncate relation " + packageRelation + " " + err.Error())
	}
}

func (t *TestCase) Run(n int) error {
	t.setup()

	for i := 0; i < n; i++ {
		go t.RunGoGet(i)
	}

	return t.Monitor(n)
}

func (t *TestCase) Monitor(n int) error {
	defer t.session.Close()
	for {
		select {
		case run := <-t.runCh:
			atomic.AddInt64(&t.runs, 1)
			if !run.success {
				return errors.New("Error in test: " + run.message)
			}
			if int(t.runs) == n {
				return nil
			}
		}
	}
}

func (t *TestCase) formatGophrUrlForIndex(index int) string {
	return fmt.Sprintf(baseGophrUrl, t.minikubeIp, repos[index], semVers[index])
}

func (t *TestCase) RunGoGet(index int) {
	gophrUrl := t.formatGophrUrlForIndex(index)
	resp, err := t.httpClient.Get(gophrUrl)
	defer resp.Body.Close()
	if err != nil {
		t.runCh <- RunInfo{success: false, message: "Failed to get gophr endpoint " + err.Error()}
		return
	}

	if resp.StatusCode != 200 {
		t.runCh <- RunInfo{success: false, message: "Failed to get gophr endpoint " + resp.Status}
		return
	}

	splitString := strings.Split(repos[index], "/")
	author, repo := splitString[0], splitString[1]
	var actualShas []string

	for i := 0; i < 5; i++ {
		time.Sleep(250 * time.Millisecond)
		actualShas = t.queryForShas(author, repo)
		if len(actualShas) != 0 {
			break
		}
	}

	if len(actualShas) == 0 {
		t.runCh <- RunInfo{success: false, message: "Failed to obtain sha for repo: " + repo + " and author " + author}
		return
	}

	if actualShas[0] != expectedShas[index] {
		t.runCh <- RunInfo{success: false, message: "Unexpected sha " + actualShas[0] + " expected " + expectedShas[index] + " for repo " + repos[index]}
	} else {
		t.runCh <- RunInfo{success: true, message: "I'm da best mayne I did it"}
	}
}

func (t *TestCase) queryForShas(author, repo string) []string {
	var actualSha string
	returnedShas := make([]string, 0)

	iter := t.session.Query("SELECT sha from " + packageRelation +
		" where author = '" + author + "' and repo = '" + repo + "'").Iter()
	for iter.Scan(&actualSha) {
		returnedShas = append(returnedShas, actualSha)
	}

	if err := iter.Close(); err != nil {
		t.runCh <- RunInfo{success: false, message: "Failed to scan gophr relation for repo " + author + " " + repo + " with error " + err.Error()}
	}

	return returnedShas
}

func TestDefaultGoGet(t *testing.T) {
	testCase := newTestCase(t)
	err := testCase.Run(4)
	if err != nil {
		assert.Fail(t, err.Error())
	}
}
