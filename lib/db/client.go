package db

import (
	"fmt"

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
	// Create the database cluster struct.
	cluster := gocql.NewCluster(c.DbAddress)
	cluster.ProtoVersion = DBProtoVersion
	if c.IsDev {
		cluster.Consistency = gocql.One
	} else {
		cluster.Consistency = gocql.Quorum
	}

	// Use the cluster to start a session.
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf(`Failed to create a database client: %v`, err)
	}

	return clientImpl{session: session}, nil
}
