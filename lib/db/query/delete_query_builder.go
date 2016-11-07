package query

import (
	"bytes"

	"github.com/gophr-pm/gophr/lib/db"
)

// DeleteQueryBuilder constructs a delete query.
type DeleteQueryBuilder struct {
	table      string
	columns    []string
	conditions []*Condition
}

// Delete starts constructing a delete query.
func Delete(columns ...string) *DeleteQueryBuilder {
	return &DeleteQueryBuilder{
		columns: columns,
	}
}

// DeleteRows starts constructing a delete query that deletes every column of
// rows that match.
func DeleteRows() *DeleteQueryBuilder {
	return &DeleteQueryBuilder{
		columns: nil,
	}
}

// From specifies the table in a delete query.
func (qb *DeleteQueryBuilder) From(table string) *DeleteQueryBuilder {
	qb.table = table
	return qb
}

// Where adds a condition to which all of the deleteed rows should adhere.
func (qb *DeleteQueryBuilder) Where(condition *Condition) *DeleteQueryBuilder {
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// And is an alias for DeleteQueryBuilder.Where(condition).
func (qb *DeleteQueryBuilder) And(condition *Condition) *DeleteQueryBuilder {
	return qb.Where(condition)
}

// Create serializes and creates the query.
func (qb *DeleteQueryBuilder) Create(q db.Queryable) db.Query {
	var (
		buffer     bytes.Buffer
		parameters []interface{}
	)

	buffer.WriteString("delete ")
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

	return q.Query(buffer.String(), parameters...)
}
