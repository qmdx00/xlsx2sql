package xlsx2sql

import (
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

// Valuer is a function that takes a string and returns a string.
type Valuer func(string) string

// xlsx ...
type xlsx struct {
	// opts is the options for the xlsx.Writer.
	opts options

	// done is the channel that is closed when the xlsx.Writer is done.
	done chan struct{}
	// channel is the channel that is used to read the rows from the spreadsheet.
	channel chan []string

	// path is the path to the spreadsheet.
	path string
	// file is the spreadsheet file.
	file *excelize.File

	// keys is the list of keys to SQL columns.
	keys []string
	// rows is the list of rows to be read from the spreadsheet.
	// Each row is a map of keys to values.
	rows []map[string]string
	// valuers is a map of keys to valuers.
	// Each valuer is a function that used to process the value
	valuers map[string]Valuer

	// addtions is the list of additional data to be added to the rows.
	addtions map[string]string
	// mapped is the map of keys to column indicies.
	mapped map[string]uint
}

// New ...
func New(path string, opts ...Option) *xlsx {
	_opts := defaultOptions
	for _, o := range opts {
		o(&_opts)
	}

	return &xlsx{
		opts:     _opts,
		channel:  make(chan []string, 10),
		done:     make(chan struct{}, 1),
		path:     path,
		file:     nil,
		keys:     make([]string, 0),
		rows:     make([]map[string]string, 0),
		valuers:  make(map[string]Valuer),
		addtions: make(map[string]string),
		mapped:   make(map[string]uint),
	}
}

// readSheet reads the spreadsheet and sends the rows to the channel.
func (x *xlsx) readSheet() (err error) {
	if x.file, err = excelize.OpenFile(x.path); err != nil {
		x.done <- struct{}{}
		return
	}
	defer x.file.Close()

	sheet := x.opts.sheet
	if x.opts.sheet == "" {
		sheet = x.file.GetSheetName(0)
	}
	rows, err := x.file.Rows(sheet)
	if err != nil {
		x.done <- struct{}{}
		return
	}

	go x.read(rows)
	go x.process()

	select {
	case <-x.done:
		return
	}
}

// Header maps the keys to the column indicies in the spreadsheet.
// The keys are stored so that they can be written in order later.
func (x *xlsx) Header(keys ...string) *xlsx {
	for i, key := range keys {
		x.mapped[key] = uint(i)
	}
	x.keys = append(x.keys, keys...)
	return x
}

// Column adds a key/value pair to the map of additional
// information. If the key already exists, Addtion replaces
// the existing value with the new value.
func (x *xlsx) Column(key string, value string) *xlsx {
	x.addtions[key] = value
	x.keys = append(x.keys, key)
	return x
}

// Mapped maps the key to the column index in the spreadsheet.
// The key is stored so that it can be written in order later.
func (x *xlsx) Mapped(key string, column uint) *xlsx {
	x.mapped[key] = column
	x.keys = append(x.keys, key)
	return x
}

// Valuer Set the valuer to be used to process the data which column mapped to the key.
func (x *xlsx) Valuer(key string, valuer Valuer) *xlsx {
	x.valuers[key] = valuer
	return x
}

// Valuers Set the valuers to be used to process the data which column mapped to the key.
func (x *xlsx) Valuers(valuers map[string]Valuer) *xlsx {
	for k, v := range valuers {
		x.valuers[k] = v
	}
	return x
}

// Build ...
func (x *xlsx) Build() (sqls []Statement, err error) {
	// read the spreadsheet
	if err = x.readSheet(); err != nil {
		return
	}
	// build the statements
	sqls = make([]Statement, 0, len(x.rows))
	switch x.opts.mode {
	case ModeSingle:
		sqls = x.singleSQL()
	case ModeBatch:
		sqls = x.batchSQL()
	}
	return
}

// ExportSQL ...
func ExportSQL(filepath string, sqls []Statement) (err error) {
	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	if err = writeSQLs(file, sqls); err != nil {
		return
	}

	log.Printf("Total: %d SQL", len(sqls))

	return
}

// ==============================================================//
// ========================= private ===========================//
// ==============================================================//

// read ...
func (x *xlsx) read(rows *excelize.Rows) {
	// skip header
	if x.opts.skipHeader {
		rows.Next()
	}
	for rows.Next() {
		if row, err := rows.Columns(); err == nil {
			x.channel <- row
		}
	}
	close(x.channel)
}

// process ...
func (x *xlsx) process() {
	for row := range x.channel {
		m := make(map[string]string)
		for _, key := range x.keys {
			// set value to empty string
			m[key] = ""
			// set value with mapped
			if column, ok := x.mapped[key]; ok {
				m[key] = row[column]
			}
			// cover value with addtions
			if val, ok := x.addtions[key]; ok {
				m[key] = val
			}
			// call valuer if exists
			if valuer, ok := x.valuers[key]; ok {
				m[key] = valuer(m[key])
			}
		}
		x.rows = append(x.rows, m)
	}
	x.done <- struct{}{}
	close(x.done)
}

// singleSQL ...
func (x *xlsx) singleSQL() (sqls []Statement) {
	sqls = make([]Statement, 0)

	var vals []string
	for _, row := range x.rows {
		vals = make([]string, 0, len(x.keys))
		for _, key := range x.keys {
			vals = append(vals, row[key])
		}
		sqls = append(sqls, &singleSQL{x.opts, x.keys, vals})
	}
	return
}

// batchSQL ...
func (x *xlsx) batchSQL() (sqls []Statement) {
	sqls = make([]Statement, 0)

	var row []string
	var vals [][]string
	for i, m := range x.rows {
		row = make([]string, 0, len(x.keys))
		for _, key := range x.keys {
			row = append(row, m[key])
		}
		if i%x.opts.batchSize == 0 {
			vals = make([][]string, 0, x.opts.batchSize)
		}
		vals = append(vals, row)
		if (i+1)%x.opts.batchSize == 0 || i == len(x.rows)-1 {
			sqls = append(sqls, &batchSQL{x.opts, x.keys, vals})
		}
	}
	return
}

// writeSQLs writes the given SQL statements to the given file, separated by newlines.
// The file is not closed after the SQL statements are written.
func writeSQLs(file *os.File, sqls []Statement) error {
	for _, sql := range sqls {
		file.WriteString(string(sql.Render()))
		file.WriteString("\n")
	}
	return nil
}
