package query

import (
	"bytes"

	"github.com/gocql/gocql"
)

// ColumnValueAssignment represents a value assignment for a specific column of
// a row.
type columnValueAssignment struct {
	column string
	value  string
}

// UpdateQueryBuilder constructs an insert query.
type UpdateQueryBuilder struct {
	valueAssignments []columnValueAssignment
	conditions       []*Condition
	table            string
}

// Update starts constructing an insert query.
func Update(table string) *UpdateQueryBuilder {
	return &UpdateQueryBuilder{
		table: table,
	}
}

// Set adds a value assignment to the update query.
func (qb *UpdateQueryBuilder) Set(column string, value string) *UpdateQueryBuilder {
	qb.valueAssignments = append(qb.valueAssignments, columnValueAssignment{
		column: column,
		value:  value,
	})
	return qb
}

// Where adds a condition to which all of the updated rows should adhere.
func (qb *UpdateQueryBuilder) Where(condition *Condition) *UpdateQueryBuilder {
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// And is an alias for UpdateQueryBuilder.Where(condition).
func (qb *UpdateQueryBuilder) And(condition *Condition) *UpdateQueryBuilder {
	return qb.Where(condition)
}

// Create serializes and creates the query.
func (qb *UpdateQueryBuilder) Create(session *gocql.Session) *gocql.Query {
	var (
		buffer     bytes.Buffer
		parameters []interface{}
	)

	buffer.WriteString("update ")
	buffer.WriteString(DBKeyspaceName)
	buffer.WriteByte('.')
	buffer.WriteString(qb.table)
	buffer.WriteString(" set ")
	for i, valueAssignment := range qb.valueAssignments {
		if i > 0 {
			buffer.WriteByte(',')
		}

		buffer.WriteString(valueAssignment.column)
		buffer.WriteString("=?")
	}
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

	return session.Query(buffer.String(), parameters)
}
