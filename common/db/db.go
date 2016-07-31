package db

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common/config"
	"github.com/skeswa/gophr/common/db/query"
	"github.com/skeswa/migrate/migrate"

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

	// Run the migrator.
	log.Println("Executing pending database migrations.")
	errs, ok := migrate.UpSync(getMigrateConnectionString(c), c.MigrationsPath)

	if !ok {
		return nil, fmt.Errorf("Database migrations failed: %v\n", errs)
	}

	return session, nil
}

// getMigrateConnectionString puts the migrate connection string together.
func getMigrateConnectionString(conf *config.Config) string {
	// Create the migrate connection string.
	buffer := bytes.Buffer{}
	buffer.WriteString("cassandra://")
	buffer.WriteString(conf.DbAddress)
	buffer.WriteByte('/')
	buffer.WriteString(query.DBKeyspaceName)
	buffer.WriteString("?protocol=")
	buffer.WriteString(strconv.Itoa(query.DBProtoVersion))

	if conf.IsDev {
		buffer.WriteString("&consistency=one&timeout=1")
	} else {
		// TODO(Skeswa): Fix it
		buffer.WriteString("&consistency=one&timeout=30")
	}

	return buffer.String()
}
