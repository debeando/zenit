package overflow

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type Column struct {
	dataType string
	unsigned bool
	current  int64
	percent  float64
	maximum  uint64
}

const (
	queryFields = `SELECT DISTINCT c.table_schema, c.table_name, c.column_name, SUBSTRING_INDEX(c.column_type, '(', 1) AS data_type
FROM information_schema.columns c
WHERE c.table_schema NOT IN ('mysql','sys','performance_schema','information_schema')
  AND c.column_type LIKE '%int%'
  AND c.column_key = 'PRI'
ORDER BY c.table_schema, c.table_name, c.column_name`
	queryMax = "SELECT COALESCE(MAX(%s), 0) AS max FROM `%s`.`%s`"
)

type MySQLOverflow struct{}

func (l *MySQLOverflow) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputMySQLOverflow", map[string]interface{}{"error": err})
		}
	}()

	for host := range config.File.Inputs.MySQL {
		if !config.File.Inputs.MySQL[host].Overflow {
			return
		}

		log.Info("InputMySQLOverflow", map[string]interface{}{
			"hostname": config.File.Inputs.MySQL[host].Hostname,
		})

		var a = metrics.Load()
		var m = mysql.GetInstance(config.File.Inputs.MySQL[host].Hostname)

		m.Connect(config.File.Inputs.MySQL[host].DSN)

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
				c.current = int64(value)

				c.Unsigned()
				c.Maximum()
				c.Percentage()

				a.Add(metrics.Metric{
					Key: "mysql_overflow",
					Tags: []metrics.Tag{
						{"hostname", config.File.Inputs.MySQL[host].Hostname},
						{"schema", rows[row]["table_schema"]},
						{"table", rows[row]["table_name"]},
						{"data_type", c.dataType},
						{"unsigned", strconv.FormatBool(c.unsigned)}},
					Values: []metrics.Value{
						{"percentage", c.percent },
					},
				})

				log.Debug("InputMySQLOverflow", map[string]interface{}{
					"hostname": config.File.Inputs.MySQL[host].Hostname,
					"schema": rows[row]["table_schema"],
					"table": rows[row]["table_name"],
					"column": rows[row]["column_name"],
					"data_type": c.dataType,
					"unsigned": c.unsigned,
					"value": value,
					"current": c.current,
					"maximum": c.maximum,
					"percent": c.percent,
				})
			}
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
	c.percent = common.Percentage(c.current, c.maximum)
}

func init() {
	inputs.Add("InputMySQLOverflow", func() inputs.Input { return &MySQLOverflow{} })
}
