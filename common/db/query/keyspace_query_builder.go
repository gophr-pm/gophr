package query

import (
	"bytes"
	"strconv"

	"github.com/gocql/gocql"
)

// KeyspaceQueryBuilder constructs an create keyspace query.
type KeyspaceQueryBuilder struct {
	durableWrites     bool
	replicationClass  string
	replicationFactor int
}

// CreateKeyspaceIfNotExists starts constructing a keyspace creation query.
func CreateKeyspaceIfNotExists() *KeyspaceQueryBuilder {
	return &KeyspaceQueryBuilder{}
}

// WithDurableWrites specifies whether writes to the keyspace should be durable.
func (qb *KeyspaceQueryBuilder) WithDurableWrites(durableWrites bool) *KeyspaceQueryBuilder {
	qb.durableWrites = durableWrites
	return qb
}

// WithReplication describes the replication of the keyspace.
func (qb *KeyspaceQueryBuilder) WithReplication(class string, factor int) *KeyspaceQueryBuilder {
	qb.replicationClass = class
	qb.replicationFactor = factor
	return qb
}

// Create serializes and creates the query.
func (qb *KeyspaceQueryBuilder) Create(session *gocql.Session) *gocql.Query {
	var buffer bytes.Buffer

	buffer.WriteString("create keyspace if not exists ")
	buffer.WriteString(DBKeyspaceName)
	buffer.WriteString(" with replication = {'class': '")
	buffer.WriteString(qb.replicationClass)
	buffer.WriteString("', 'replication_factor': '")
	buffer.WriteString(strconv.Itoa(qb.replicationFactor))
	buffer.WriteString("'} and durable_writes = ")
	buffer.WriteString(strconv.FormatBool(qb.durableWrites))

	return session.Query(buffer.String())
}
