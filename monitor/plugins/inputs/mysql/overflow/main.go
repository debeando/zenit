package overflow

import (
	"fmt"
	"strconv"
	"strings"

	"zenit/config"
	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/lists/metrics"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/math"
	"github.com/debeando/go-common/mysql"
)

type Column struct {
	dataType string
	unsigned bool
	current  int64
	percent  float64
	maximum  uint64
}

const (
	queryFields = `SELECT DISTINCT
	c.table_schema,
	c.table_name,
	c.column_name,
	SUBSTRING_INDEX(c.column_type, '(', 1) AS data_type
FROM information_schema.columns c
WHERE c.table_schema NOT IN ('mysql','sys','performance_schema','information_schema')
  AND c.column_type LIKE '%int%'
  AND c.column_key = 'PRI'
ORDER BY c.table_schema, c.table_name, c.column_name`
	queryMax = "SELECT COALESCE(MAX(%s), 0) AS max FROM `%s`.`%s`"
)

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	for host := range cnf.Inputs.MySQL {
		log.DebugWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
			"enable":   cnf.Inputs.MySQL[host].Enable,
			"overflow": cnf.Inputs.MySQL[host].Overflow,
		})

		if !cnf.Inputs.MySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.MySQL[host].Overflow {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
		})

		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		err := m.Connect()
		if err != nil {
			continue
		}

		rows, _ := m.Query(queryFields)
		for row := range rows {
			log.DebugWithFields(name, log.Fields{
				"hostname":  cnf.Inputs.MySQL[host].Hostname,
				"database":  rows[row]["TABLE_SCHEMA"],
				"table":     rows[row]["TABLE_NAME"],
				"column":    rows[row]["COLUMN_NAME"],
				"data_type": rows[row]["data_type"],
			})

			max, _ := m.Query(
				fmt.Sprintf(
					queryMax,
					rows[row]["COLUMN_NAME"],
					rows[row]["TABLE_SCHEMA"],
					rows[row]["TABLE_NAME"],
				),
			)

			if value, ok := mysql.ParseValue(max[0]["max"]); ok {
				var c Column
				c.dataType = rows[row]["data_type"]
				c.current = int64(value)

				c.Unsigned()
				c.Maximum()
				c.Percentage()

				mtc.Add(metrics.Metric{
					Key: "mysql_overflow",
					Tags: []metrics.Tag{
						{Name: "hostname", Value: cnf.Inputs.MySQL[host].Hostname},
						{Name: "schema", Value: rows[row]["TABLE_SCHEMA"]},
						{Name: "table", Value: rows[row]["TABLE_NAME"]},
						{Name: "data_type", Value: c.dataType},
						{Name: "unsigned", Value: strconv.FormatBool(c.unsigned)}},
					Values: []metrics.Value{
						{Key: "percentage", Value: c.percent},
					},
				})

				log.DebugWithFields(name, log.Fields{
					"hostname":  cnf.Inputs.MySQL[host].Hostname,
					"schema":    rows[row]["TABLE_SCHEMA"],
					"table":     rows[row]["TABLE_NAME"],
					"column":    rows[row]["COLUMN_NAME"],
					"data_type": c.dataType,
					"unsigned":  c.unsigned,
					"value":     value,
					"current":   c.current,
					"maximum":   c.maximum,
					"percent":   c.percent,
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
	c.percent = math.Percentage(c.current, c.maximum)
}

func init() {
	inputs.Add("InputMySQLOverflow", func() inputs.Input { return &Plugin{} })
}
