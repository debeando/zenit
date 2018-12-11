// Overflow collect the max value of Primary Key on table and verify the limit which Data Type.

package overflow

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

// Column is a struct to save result of query.
type Column struct {
	schema   string
	table    string
	column   string
	dataType string
	unsigned bool
	current  float64
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

type MySQLOverflow struct {}

func (l *MySQLOverflow) Collect() {
	if ! config.File.MySQL.Inputs.Overflow {
		return
	}

	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Error("MySQL:Overflow - Impossible to connect: " + err.Error())
	}

	rows, err := conn.Query(querySQLColumns)
	defer rows.Close()
	if err != nil {
		log.Error("MySQL:Overflow - Impossible to execute query: " + err.Error())
	}

	var a = metrics.Load()

	for rows.Next() {
		var c Column

		rows.Scan(
			&c.schema,
			&c.table,
			&c.column,
			&c.dataType)

		err = conn.QueryRow(fmt.Sprintf(querySQLMaxInt, c.column, c.schema, c.table)).Scan(&c.current)
		if err != nil {
			log.Error("MySQL:Overflow - Impossible to execute query: " + err.Error())
		}

		c.Unsigned()
		c.Maximum()
		c.Percentage()

		a.Add(metrics.Metric{
			Key: "zenit_mysql_stats_overflow",
			Tags: []metrics.Tag{
				{"schema", c.schema},
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
	c.percent = common.Percentage(c.current, float64(c.maximum))
}

func init() {
	loader.Add("MySQLOverflow", func() loader.Plugin { return &MySQLOverflow{} })
}
