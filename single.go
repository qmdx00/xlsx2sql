package xlsx2sql

import "strings"

// singleSQL is a type that implements the Statement interface.
var _ Statement = (*singleSQL)(nil)

// singleSQL is used to create a single SQL statement.
// A single SQL statement is a single SQL statement that contains a single
// row of data.
type singleSQL struct {
	// opts is a variable that contains the options.
	opts options
	// keys is a slice of strings that contains the column names.
	keys []string
	// vals is a slice of strings that contains the values.
	vals []string
}

// mode ...
func (*singleSQL) mode() string {
	return ModeSingle
}

// Render returns a string of SQL that can be used to insert the values
// into the database. The string is formatted as an INSERT statement.
func (s *singleSQL) Render() string {
	var sb strings.Builder
	// Write the beginning of the statement.
	sb.WriteString("INSERT INTO `")
	sb.WriteString(s.opts.table)
	sb.WriteString("` (")

	// Write the keys.
	for i, key := range s.keys {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("`")
		sb.WriteString(key)
		sb.WriteString("`")
	}
	sb.WriteString(") VALUES (")

	// Write the values.
	for i, val := range s.vals {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("'")
		sb.WriteString(val)
		sb.WriteString("'")
	}
	// Write the end of the statement.
	sb.WriteString(");")
	return sb.String()
}
