package xlsx2sql

import "strings"

// batchSQL is a type that implements the Statement interface.
var _ Statement = (*batchSQL)(nil)

// batchSQL is used to create a batch SQL statement.
// A batch SQL statement is a single SQL statement that contains multiple
// rows of data.
type batchSQL struct {
	// opts is a variable that contains the options.
	opts options
	// keys is a slice of strings that contains the column names.
	keys []string
	// vals is a slice of strings that contains the values.
	vals [][]string
}

// mode ...
func (s *batchSQL) mode() string {
	return ModeBatch
}

// Render returns a string of SQL that can be used to insert the values
// into the database. The string is formatted as an INSERT statement.
func (s *batchSQL) Render() string {
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
	sb.WriteString(") VALUES ")
	// Write the values.
	for i, val := range s.vals {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("(")
		for j, v := range val {
			if j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString("'")
			sb.WriteString(v)
			sb.WriteString("'")
		}
		sb.WriteString(")")
	}
	// Write the end of the statement.
	sb.WriteString(";")
	return sb.String()
}
