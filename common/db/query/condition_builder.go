package query

import "bytes"

// Condition is a query filter used in db queries.
type Condition struct {
	expression   string
	parameter    interface{}
	hasParameter bool
}

// ColumnConditionBuilder helps to buidl conditions.
type ColumnConditionBuilder struct {
	column string
}

// Column creates a condition builder.
func Column(column string) *ColumnConditionBuilder {
	return &ColumnConditionBuilder{
		column: column,
	}
}

// Equals creates an equivalence condition.
func (cb *ColumnConditionBuilder) Equals(value interface{}) *Condition {
	return &Condition{
		expression:   cb.column + "=?",
		parameter:    value,
		hasParameter: true,
	}
}

// IndexConditionBuilder helps to build conditions.
type IndexConditionBuilder struct {
	index string
}

// Index creates a condition builder.
func Index(index string) *IndexConditionBuilder {
	return &IndexConditionBuilder{
		index: index,
	}
}

// Query creates an index condition.
func (cb *IndexConditionBuilder) Query(query string) *Condition {
	var buffer bytes.Buffer

	buffer.WriteString("expr(")
	buffer.WriteString(cb.index)
	buffer.WriteString(",'")
	buffer.WriteString(query)
	buffer.WriteString("')")

	return &Condition{
		expression:   buffer.String(),
		hasParameter: false,
	}
}
