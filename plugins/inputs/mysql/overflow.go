package mysql

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

// Column is a struct to save result of query.
type Column struct {
	schema   string
	table    string
	column   string
	dataType string
	unsigned bool
	current  uint64
	percent  float64
	maximum  uint64
}

const (
	querySQLColumns = `SELECT table_schema, table_name, column_name, SUBSTRING_INDEX(column_type, '(', 1) AS column_type
FROM information_schema.columns
WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema','percona')
  AND column_type LIKE '%int%'
  AND column_key = 'PRI'`
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

		rows.Scan(
			&c.schema,
			&c.table,
			&c.column,
			&c.dataType)

		err = conn.QueryRow(fmt.Sprintf(querySQLMaxInt, c.column, c.schema, c.table)).Scan(&c.current)
		if err != nil {
			log.Printf("E! - MySQL:Overflow - Impossible to execute query: %s\n", err)
		}

		c.Unsigned()
		c.Maximum()
		c.Percentage()

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

func (c *Column) Unsigned() {
	c.unsigned = strings.Contains(c.dataType, "unsigned")
}

func (c *Column) Maximum() {
	if c.unsigned == true {
		c.maximum = mysql.MaximumValueUnsigned(c.dataType)
	} else {
		c.maximum = mysql.MaximumValueSigned(c.dataType)
	}
}

func (c *Column) Percentage() {
	c.percent = common.Percentage(c.current, c.maximum)
}
