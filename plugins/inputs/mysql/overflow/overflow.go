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

type Column struct {
	dataType string
	unsigned bool
	current  uint64
	percent  float64
	maximum  uint64
}

const (
	queryFields = `SELECT DISTINCT c.table_schema, c.table_name, c.column_name, SUBSTRING_INDEX(c.column_type, '(', 1) AS data_type
FROM information_schema.columns c
WHERE c.table_schema NOT IN ('mysql','sys','performance_schema','information_schema','percona')
  AND c.column_type LIKE '%int%'
  AND c.column_key = 'PRI'
ORDER BY c.table_schema, c.table_name, c.column_name`
	queryMax = "SELECT COALESCE(MAX(%s), 0) AS max FROM %s.%s"
)

type MySQLOverflow struct {}

func (l *MySQLOverflow) Collect() {
	if ! config.File.MySQL.Inputs.Overflow {
		return
	}

	var a = metrics.Load()
	var m = mysql.GetInstance("mysql")
	m.Connect(config.File.MySQL.DSN)

	rows := m.Query(queryFields)

	for row := range rows {
		max := m.Query(
			fmt.Sprintf(
				queryMax,
				rows[row]["column_name"],
				rows[row]["table_schema"],
				rows[row]["table_name"],
			),
		)

		if value, ok := mysql.ParseValue(max[0]["max"]); ok {
			var c Column
			c.dataType = rows[row]["data_type"]
			c.current = value

			c.Unsigned()
			c.Maximum()
			c.Percentage()

			a.Add(metrics.Metric{
				Key: "zenit_mysql_overflow",
				Tags: []metrics.Tag{
					{"schema", rows[row]["table_schema"],},
					{"table", rows[row]["table_name"]},
					{"type", "overflow"},
					{"data_type", c.dataType},
					{"unsigned", strconv.FormatBool(c.unsigned)}},
				Values: c.percent,
			})

			log.Debug(
				fmt.Sprintf("Plugin - InputMySQLOverflow - %s.%s.%s(%s,%t)=%d [(%d/%d)*100=%.2f%%]",
					rows[row]["table_schema"],
					rows[row]["table_name"],
					rows[row]["column_name"],
					c.dataType,
					c.unsigned,
					value,
					c.current,
					c.maximum,
					c.percent,
				),
			)
		}
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
	c.percent = common.Percentage(float64(c.current), float64(c.maximum))
}

func init() {
	loader.Add("InputMySQLOverflow", func() loader.Plugin { return &MySQLOverflow{} })
}
