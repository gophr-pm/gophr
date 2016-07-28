package query

import (
	"bytes"

	"github.com/gocql/gocql"
)

// InsertQueryBuilder constructs an insert query.
type InsertQueryBuilder struct {
	columns []string
	values  []interface{}
	table   string
}

// InsertInto starts constructing an insert query.
func InsertInto(table string) *InsertQueryBuilder {
	return &InsertQueryBuilder{
		table: table,
	}
}

// Value adds a value mapping to the insert query.
func (qb *InsertQueryBuilder) Value(column string, value interface{}) *InsertQueryBuilder {
	qb.columns = append(qb.columns, column)
	qb.values = append(qb.values, value)
	return qb
}

// Create serializes and creates the query.
func (qb *InsertQueryBuilder) Create(session *gocql.Session) *gocql.Query {
	var buffer bytes.Buffer

	buffer.WriteString("insert into ")
	buffer.WriteString(DBKeyspaceName)
	buffer.WriteByte('.')
	buffer.WriteString(qb.table)
	buffer.WriteString(" (")
	for i, column := range qb.columns {
		if i > 0 {
			buffer.WriteByte(',')
		}

		buffer.WriteString(column)
	}
	buffer.WriteString(") values (")
	for i := range qb.columns {
		if i > 0 {
			buffer.WriteByte(',')
		}

		buffer.WriteByte('?')
	}
	buffer.WriteByte(')')

	return session.Query(buffer.String(), qb.values...)
}
