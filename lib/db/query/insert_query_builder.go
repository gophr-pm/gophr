package query

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

// InsertionValueType is the type of value in the InsertQueryBuilder#Value(...).
type InsertionValueType int

const (
	// LiteralValue is a value that need not be parameterized.
	LiteralValue InsertionValueType = iota
	// ParameterizedValue is a value that should be parameterized.
	ParameterizedValue
)

// InsertQueryBuilder constructs an insert query.
type InsertQueryBuilder struct {
	ifNotExists bool
	valueTypes  []InsertionValueType
	columns     []string
	values      []interface{}
	table       string
	ttl         time.Duration
}

// InsertInto starts constructing an insert query.
func InsertInto(table string) *InsertQueryBuilder {
	return &InsertQueryBuilder{
		table: table,
	}
}

// Value adds a value mapping to the insert query.
func (qb *InsertQueryBuilder) Value(column string, value interface{}, types ...InsertionValueType) *InsertQueryBuilder {
	var valueType InsertionValueType
	switch len(types) {
	case 0:
		// ParameterizedValue is the default.
		valueType = ParameterizedValue
	case 1:
		valueType = types[0]
	default:
		panic("invalid insertion value type")
	}

	qb.valueTypes = append(qb.valueTypes, valueType)
	qb.columns = append(qb.columns, column)
	qb.values = append(qb.values, value)
	return qb
}

// UsingTTL adds a TTL to the query.
func (qb *InsertQueryBuilder) UsingTTL(ttl time.Duration) *InsertQueryBuilder {
	qb.ttl = ttl
	return qb
}

// IfNotExists makes the query only perform a write if unique.
func (qb *InsertQueryBuilder) IfNotExists() *InsertQueryBuilder {
	qb.ifNotExists = true
	return qb
}

// Create serializes and creates the query.
func (qb *InsertQueryBuilder) Create(q Queryable) *gocql.Query {
	var (
		buffer bytes.Buffer
		params []interface{}
	)

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

		switch qb.valueTypes[i] {
		case LiteralValue:
			buffer.WriteString(fmt.Sprintf("%v", qb.values[i]))
		case ParameterizedValue:
			buffer.WriteByte('?')
			params = append(params, qb.values[i])
		default:
			panic(fmt.Sprintf("invalid value type %d encountered", qb.valueTypes[i]))
		}

	}
	buffer.WriteByte(')')
	if qb.ifNotExists {
		buffer.WriteString(" if not exists")
	}
	if qb.ttl > 0 {
		buffer.WriteString(" using TTL ")
		buffer.WriteString(strconv.FormatUint(uint64(qb.ttl/time.Microsecond), 10))
	}

	return q.Query(buffer.String(), params...)
}
