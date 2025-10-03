package merger

import "luamerge/internal/parser"

// Result represents the result of a merge operation.
// Contains the table name and the table data after merging.
type Result struct {
	TableName string
	Table     *parser.Table
}
