package db

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/skeswa/gophr/common/db/query"
	"github.com/skeswa/migrate/migrate"
)

// Migrate runs all pending migrations on the database addressed in conf.
func Migrate(dbAddress string, migrationsPath string) error {
	// Create the migrate connection string.
	buffer := bytes.Buffer{}
	buffer.WriteString("cassandra://")
	buffer.WriteString(dbAddress)
	buffer.WriteByte('/')
	buffer.WriteString(query.DBKeyspaceName)
	buffer.WriteString("?protocol=")
	buffer.WriteString(strconv.Itoa(query.DBProtoVersion))
	buffer.WriteString("&consistency=one&timeout=10")

	if errs, ok := migrate.UpSync(buffer.String(), migrationsPath); !ok {
		return fmt.Errorf("Database migrations failed: %v.", errs)
	}

	return nil
}
