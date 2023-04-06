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
	sql := "INSERT INTO `" + s.opts.table + "` ("

	keys := make([]string, 0, len(s.keys))
	for _, key := range s.keys {
		keys = append(keys, "`"+key+"`")
	}
	sql += strings.Join(keys, ", ")
	sql += ") VALUES "

	rows := make([]string, 0, len(s.vals))
	var row string
	for _, val := range s.vals {
		row = "("
		for j, v := range val {
			val[j] = "`" + v + "`"
		}
		row += strings.Join(val, ", ")
		row += ")"
		rows = append(rows, row)
	}
	sql += strings.Join(rows, ", ")
	sql += ";"
	return sql
}
