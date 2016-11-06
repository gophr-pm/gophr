package query

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/gophr-pm/gophr/lib/db"
)

// ColumnValueAssignment represents a value assignment for a specific column of
// a row.
type columnValueAssignment struct {
	column        string
	value         interface{}
	parameterized bool
}

// UpdateQueryBuilder constructs an insert query.
type UpdateQueryBuilder struct {
	valueAssignments []columnValueAssignment
	conditions       []*Condition
	ifExists         bool
	table            string
}

// Update starts constructing an insert query.
func Update(table string) *UpdateQueryBuilder {
	return &UpdateQueryBuilder{
		table: table,
	}
}

// Set adds a value assignment to the update query.
func (qb *UpdateQueryBuilder) Set(
	column string,
	value interface{},
) *UpdateQueryBuilder {
	qb.valueAssignments = append(qb.valueAssignments, columnValueAssignment{
		column:        column,
		value:         value,
		parameterized: true,
	})
	return qb
}

// Increment increases the value of a counter by a specified amount.
func (qb *UpdateQueryBuilder) Increment(
	column string,
	amount int,
) *UpdateQueryBuilder {
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

// IfExists signals that this query should only be applied to existing rows.
func (qb *UpdateQueryBuilder) IfExists() *UpdateQueryBuilder {
	qb.ifExists = true
	return qb
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
			buffer.WriteString(fmt.Sprint(valueAssignment.value))
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
	if qb.ifExists {
		buffer.WriteString(" if exists")
	}

	return buffer.String(), parameters
}

// Create serializes and creates the query.
func (qb *UpdateQueryBuilder) Create(q db.Queryable) db.Query {
	text, params := qb.compose()
	return q.Query(text, params...)
}

// AppendTo serializes and creates the query with no return.
func (qb *UpdateQueryBuilder) AppendTo(b db.Batch) {
	text, params := qb.compose()
	b.Query(text, params...)
}
