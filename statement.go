package xlsx2sql

const (
	// ModeSingle is the mode for a single SQL statement.
	ModeSingle = "single"
	// ModeBatch is the mode for a batch SQL statement.
	ModeBatch = "batch"
)

// Statement is an interface that defines the methods that are required for a type
// to be a SQL statement.
type Statement interface {
	// The mode method returns the mode of the SQL statement.
	// The mode can be either "single" or "batch".
	// The mode is used to determine which type of SQL statement to create.
	mode() string

	// The Render method returns a string that can be used to insert the values
	// into the database. The string is formatted as an INSERT statement.
	Render() string
}
