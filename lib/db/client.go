package db

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/lib/config"
)

const (
	// DBProtoVersion is the cassandra protocol version used by gophr.
	DBProtoVersion = 4
)

// Client interfaces with the cassandra database.
type Client interface {
	BatchingQueryable

	// Close closes all connections. The client is unusable after this operation.
	Close()
}

// NewClient creates a new database client.
func NewClient(c *config.Config) (Client, error) {
	// Get the database nodes from k8s.
	dbNodeIPs, err := getDBNodes()
	if err != nil {
		return nil, fmt.Errorf(`Failed to create a new database client: %v.`, err)
	}

	// Create the cluster from the IPs.
	cluster := gocql.NewCluster(dbNodeIPs...)
	// Some queries can be very expensive. Need more time to respond.
	cluster.Timeout = 2 * time.Second
	// Makes for better performance on contrained memory.
	cluster.PageSize = 2000
	// Use latest versions of Cassandra.
	cluster.ProtoVersion = DBProtoVersion

	if c.IsDev {
		// In development, keep the load light.
		cluster.Consistency = gocql.One
	} else {
		cluster.Consistency = gocql.Quorum

		// In production, it is important that queries are tried multiple times
		// since the requests against the database are generally more important.
		cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 2}
	}

	// Use the cluster to start a session.
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf(`Failed to create a database client: %v`, err)
	}

	return clientImpl{session: session}, nil
}
