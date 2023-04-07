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
	sql := "INSERT INTO `" + s.opts.table + "` ("

	keys := make([]string, 0, len(s.keys))
	for _, key := range s.keys {
		keys = append(keys, "`"+key+"`")
	}
	sql += strings.Join(keys, ", ")
	sql += ") VALUES ("

	vals := make([]string, 0, len(s.vals))
	for _, val := range s.vals {
		vals = append(vals, "'"+val+"'")
	}
	sql += strings.Join(vals, ", ")
	sql += ");"
	return sql
}
