package db

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/skeswa/gophr/common/db/query"
	"github.com/skeswa/migrate/migrate"
)

// Migrate runs all pending migrations on the database addressed in conf.
func Migrate(isDev bool, dbAddress string, migrationsPath string) error {
	// Create the migrate connection string.
	buffer := bytes.Buffer{}
	buffer.WriteString("cassandra://")
	buffer.WriteString(dbAddress)
	buffer.WriteByte('/')
	buffer.WriteString(query.DBKeyspaceName)
	buffer.WriteString("?protocol=")
	buffer.WriteString(strconv.Itoa(query.DBProtoVersion))

	// NB: consistency is now always "all". This is due to the fact that every
	// environment now has at least two nodes.
	if isDev {
		buffer.WriteString("&consistency=all&timeout=10")
	} else {
		buffer.WriteString("&consistency=all&timeout=30")
	}

	if errs, ok := migrate.UpSync(buffer.String(), migrationsPath); !ok {
		return fmt.Errorf("Database migrations failed: %v.", errs)
	}

	return nil
}
