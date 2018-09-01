package mysql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

type Column struct {
	schema    string
	table     string
	column    string
	data_type string
	unsigned  bool
	current   uint64
	percent   float64
}

const (
	dt_tinyint        uint8  = 127
	dt_smallint       uint16 = 32767
	dt_mediumint      uint32 = 8388607
	dt_int            uint32 = 2147483647
	dt_bigint         uint64 = 9223372036854775807
	dt_us_tinyint     uint8  = 255
	dt_us_smallint    uint16 = 65535
	dt_us_mediumint   uint32 = 16777215
	dt_us_int         uint32 = 4294967295
	dt_us_bigint      uint64 = 18446744073709551615
	QUERY_SQL_COLUMNS        = `
SELECT table_schema, table_name, column_name, column_type
FROM information_schema.columns
WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema','percona')
  AND column_type LIKE '%int%'
  AND column_key = 'PRI'
`
	QUERY_SQL_MAX_INT = "SELECT COALESCE(MAX(%s), 0) FROM %s.%s"
)

func Overflow() {
	conn, err := mysql.Connect(config.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	rows, err := conn.Query(QUERY_SQL_COLUMNS)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	var a = accumulator.Load()

	for rows.Next() {
		var c Column
		var m uint64

		rows.Scan(
			&c.schema,
			&c.table,
			&c.column,
			&c.data_type)

		err = conn.QueryRow(fmt.Sprintf(QUERY_SQL_MAX_INT, c.column, c.schema, c.table)).Scan(&m)
		if err != nil {
			panic(err)
		}

		c.unsigned = strings.Contains(c.data_type, "unsigned")
		c.data_type = c.data_type[0:strings.Index(c.data_type, "(")]
		c.current = m

		if c.unsigned == true {
			switch c.data_type {
			case "tinyint":
				c.percent = (float64(c.current) / float64(dt_us_tinyint)) * 100
			case "smallint":
				c.percent = (float64(c.current) / float64(dt_us_smallint)) * 100
			case "mediumint":
				c.percent = (float64(c.current) / float64(dt_us_mediumint)) * 100
			case "int":
				c.percent = (float64(c.current) / float64(dt_us_int)) * 100
			case "bigint":
				c.percent = (float64(c.current) / float64(dt_us_bigint)) * 100
			}
		} else {
			switch c.data_type {
			case "tinyint":
				c.percent = (float64(c.current) / float64(dt_tinyint)) * 100
			case "smallint":
				c.percent = (float64(c.current) / float64(dt_smallint)) * 100
			case "mediumint":
				c.percent = (float64(c.current) / float64(dt_mediumint)) * 100
			case "int":
				c.percent = (float64(c.current) / float64(dt_int)) * 100
			case "bigint":
				c.percent = (float64(c.current) / float64(dt_bigint)) * 100
			}
		}

		a.AddItem(accumulator.Metric{
			Key: "mysql_stats_overflow",
			Tags: []accumulator.Tag{{"schema", c.schema},
				{"table", c.table},
				{"type", "overflow"},
				{"data_type", c.data_type},
				{"unsigned", strconv.FormatBool(c.unsigned)}},
			Values: c.percent,
		})
	}
}
