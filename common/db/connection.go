package db

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common/config"
	"github.com/skeswa/gophr/common/db/query"

	// Load the migrate cassandra driver.
	_ "github.com/skeswa/migrate/driver/cassandra"
)

// OpenConnection starts a database session.
func OpenConnection(c *config.Config) (*gocql.Session, error) {
	// Create the database cluster struct.
	log.Println("Creating database session.")
	cluster := gocql.NewCluster(c.DbAddress)
	cluster.ProtoVersion = query.DBProtoVersion
	if c.IsDev {
		cluster.Consistency = gocql.One
	} else {
		cluster.Consistency = gocql.Quorum
	}

	// Use the cluster to start a session.
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	// Make sure that the keyspace exists.
	log.Println("Asserting the existence of the keyspace.")
	err = query.CreateKeyspace().
		WithReplication("SimpleStrategy", 3).
		WithDurableWrites(true).
		Create(session).
		Exec()
	if err != nil {
		return nil, err
	}

	return session, nil
}
