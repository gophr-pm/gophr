package query

import (
	"bytes"
	"strconv"

	"github.com/gocql/gocql"
)

// SelectQueryBuilder constructs a select query.
type SelectQueryBuilder struct {
	columns    []string
	table      string
	conditions []*Condition
	limit      *int
}

// Select starts constructing a select query.
func Select(columns ...string) *SelectQueryBuilder {
	return &SelectQueryBuilder{
		columns: columns,
	}
}

// From specifies the table in a select query.
func (qb *SelectQueryBuilder) From(table string) *SelectQueryBuilder {
	qb.table = table
	return qb
}

// Where adds a condition to which all of the selected rows should adhere.
func (qb *SelectQueryBuilder) Where(condition *Condition) *SelectQueryBuilder {
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// And is an alias for SelectQueryBuilder.Where(condition).
func (qb *SelectQueryBuilder) And(condition *Condition) *SelectQueryBuilder {
	return qb.Where(condition)
}

// Limit specifies the maximum number of results to fetch.
func (qb *SelectQueryBuilder) Limit(limit int) *SelectQueryBuilder {
	limitClone := limit
	qb.limit = &limitClone
	return qb
}

// Create serializes and creates the query.
func (qb *SelectQueryBuilder) Create(session *gocql.Session) *gocql.Query {
	var (
		buffer     bytes.Buffer
		parameters []interface{}
	)

	buffer.WriteString("select ")
	for i, column := range qb.columns {
		if i > 0 {
			buffer.WriteByte(',')
		}

		buffer.WriteString(column)
	}
	buffer.WriteString(" from ")
	buffer.WriteString(DBKeyspaceName)
	buffer.WriteByte('.')
	buffer.WriteString(qb.table)
	if qb.conditions != nil {
		buffer.WriteString(" where ")
		for i, cond := range qb.conditions {
			if i > 0 {
				buffer.WriteString(" and ")
			}

			if cond.hasParameter {
				parameters = append(parameters, cond.parameter)
			}

			buffer.WriteString(cond.expression)
		}
	}
	if qb.limit != nil {
		buffer.WriteString(" limit ")
		buffer.WriteString(strconv.Itoa(*qb.limit))
	}

	return session.Query(buffer.String(), parameters)
}
