package query

import (
	"bytes"
	"strconv"

	"github.com/gocql/gocql"
)

// ColumnValueAssignment represents a value assignment for a specific column of
// a row.
type columnValueAssignment struct {
	column        string
	value         string
	parameterized bool
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
		column:        column,
		value:         value,
		parameterized: true,
	})
	return qb
}

// Increment increases the value of a counter by a specified amount.
func (qb *UpdateQueryBuilder) Increment(column string, amount int) *UpdateQueryBuilder {
	qb.valueAssignments = append(qb.valueAssignments, columnValueAssignment{
		column:        column,
		value:         (column + "+" + strconv.Itoa(amount)),
		parameterized: false,
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

// compose composes the text and parameters for this query.
func (qb *UpdateQueryBuilder) compose() (string, []interface{}) {
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
		buffer.WriteByte('=')
		if valueAssignment.parameterized {
			buffer.WriteByte('?')
			parameters = append(parameters, valueAssignment.value)
		} else {
			buffer.WriteString(valueAssignment.value)
		}
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

	return buffer.String(), parameters
}

// Create serializes and creates the query.
func (qb *UpdateQueryBuilder) Create(q Queryable) *gocql.Query {
	text, params := qb.compose()
	return q.Query(text, params...)
}

// CreateVoid serializes and creates the query with no return.
func (qb *UpdateQueryBuilder) CreateVoid(q VoidQueryable) {
	text, params := qb.compose()
	q.Query(text, params...)
}
