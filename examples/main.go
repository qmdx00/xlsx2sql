package main

import (
	"fmt"
	"strconv"

	"github.com/qmdx00/xlsx2sql"
)

func main() {
	sqls, err := xlsx2sql.New("../data/test.xlsx", xlsx2sql.SQLMode(xlsx2sql.ModeBatch), xlsx2sql.BatchSize(10), xlsx2sql.SkipHeader(true), xlsx2sql.TableName("test"), xlsx2sql.SheetName("Sheet1")).
		Mapped("idx", 0).
		Header("id", "name").
		Column("foo", "bar").
		Valuer("idx", func(v string) string {
			id, _ := strconv.Atoi(v)
			return strconv.Itoa(id - 1)
		}).
		Build()

	if err != nil {
		panic(err)
	}

	for _, sql := range sqls {
		fmt.Println(sql.Render())
	}

	xlsx2sql.ExportSQL("../data/test.sql", sqls)
}
