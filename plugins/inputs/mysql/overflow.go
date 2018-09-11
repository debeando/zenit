package mysql

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

// Column is a struct to save result of query.
type Column struct {
	schema    string
	table     string
	column    string
	dataType  string
	unsigned  bool
	current   uint64
	percent   float64
}

const (
	dtTinyInt       uint8  = 127
	dtSmallInt      uint16 = 32767
	dtMediumInt     uint32 = 8388607
	dtInt           uint32 = 2147483647
	dtBigInt        uint64 = 9223372036854775807
	dtUSTinyInt     uint8  = 255
	dtUSSmallInt    uint16 = 65535
	dtUSMediumInt   uint32 = 16777215
	dtUSInt         uint32 = 4294967295
	dtUSBigInt      uint64 = 18446744073709551615
	querySQLColumns        = `
SELECT table_schema, table_name, column_name, column_type
FROM information_schema.columns
WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema','percona')
  AND column_type LIKE '%int%'
  AND column_key = 'PRI'
`
	querySQLMaxInt = "SELECT COALESCE(MAX(%s), 0) FROM %s.%s"
)

// Overflow collect the max value of Primary Key on table and verify the limit which Data Type.
func Overflow() {
	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Printf("E! - MySQL:Overflow - Impossible to connect: %s\n", err)
	}

	rows, err := conn.Query(querySQLColumns)
	defer rows.Close()
	if err != nil {
		log.Printf("E! - MySQL:Overflow - Impossible to execute query: %s\n", err)
	}

	var a = accumulator.Load()

	for rows.Next() {
		var c Column
		var m uint64

		rows.Scan(
			&c.schema,
			&c.table,
			&c.column,
			&c.dataType)

		err = conn.QueryRow(fmt.Sprintf(querySQLMaxInt, c.column, c.schema, c.table)).Scan(&m)
		if err != nil {
			panic(err)
		}

		c.unsigned = strings.Contains(c.dataType, "unsigned")
		c.dataType = c.dataType[0:strings.Index(c.dataType, "(")]
		c.current = m

		if c.unsigned == true {
			switch c.dataType {
			case "tinyint":
				c.percent = (float64(c.current) / float64(dtUSTinyInt)) * 100
			case "smallint":
				c.percent = (float64(c.current) / float64(dtUSSmallInt)) * 100
			case "mediumint":
				c.percent = (float64(c.current) / float64(dtUSMediumInt)) * 100
			case "int":
				c.percent = (float64(c.current) / float64(dtUSInt)) * 100
			case "bigint":
				c.percent = (float64(c.current) / float64(dtUSBigInt)) * 100
			}
		} else {
			switch c.dataType {
			case "tinyint":
				c.percent = (float64(c.current) / float64(dtTinyInt)) * 100
			case "smallint":
				c.percent = (float64(c.current) / float64(dtSmallInt)) * 100
			case "mediumint":
				c.percent = (float64(c.current) / float64(dtMediumInt)) * 100
			case "int":
				c.percent = (float64(c.current) / float64(dtInt)) * 100
			case "bigint":
				c.percent = (float64(c.current) / float64(dtBigInt)) * 100
			}
		}

		a.Add(accumulator.Metric{
			Key: "mysql_stats_overflow",
			Tags: []accumulator.Tag{{"schema", c.schema},
				{"table", c.table},
				{"type", "overflow"},
				{"data_type", c.dataType},
				{"unsigned", strconv.FormatBool(c.unsigned)}},
			Values: c.percent,
		})
	}
}
