package xlsx2sql

import "github.com/xuri/excelize/v2"

// Option is a function that sets an option.
type Option func(*options)

// options ...
type options struct {
	excelizeOpts []excelize.Options
	mode         string
	table        string
	skipHeader   bool
	batchSize    int
	sheet        string
}

// defaultOptions is a variable that contains a set of default options.
// The default options are used when the user does not specify any options
// when calling a function.
// The default options are: "single" mode, batch size of 10, and skipHeader
// set to true.
var defaultOptions = options{
	excelizeOpts: make([]excelize.Options, 0),
	mode:         "single",
	batchSize:    10,
	skipHeader:   true,
}

// SQLMode is an option that sets the SQL mode.
func SQLMode(mode string) Option {
	return func(o *options) {
		o.mode = mode
	}
}

// BatchSize is an option that sets the batch size for the SQL statement.
func BatchSize(size int) Option {
	return func(o *options) {
		o.batchSize = size
	}
}

// SkipHeader is an option that sets whether the xlsx header should be skipped
func SkipHeader(skip bool) Option {
	return func(o *options) {
		o.skipHeader = skip
	}
}

// TableName is an option that sets the table name for the SQL statement.
func TableName(table string) Option {
	return func(o *options) {
		o.table = table
	}
}

// Sheet is an option that sets the sheet name for the SQL statement.
func SheetName(sheet string) Option {
	return func(o *options) {
		o.sheet = sheet
	}
}

// ExcelizeOptions is an option that sets the options for the excelize package.
func ExcelizeOptions(opts ...excelize.Options) Option {
	return func(o *options) {
		o.excelizeOpts = opts
	}
}
