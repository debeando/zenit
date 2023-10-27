package overflow

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

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
	SQLOverflow = `SELECT DISTINCT
	c.table_schema,
	c.table_name,
	c.column_name,
	SUBSTRING_INDEX(c.column_type, '(', 1) AS data_type
FROM information_schema.columns c
WHERE c.table_schema NOT IN ('mysql','sys','performance_schema','information_schema')
  AND c.column_type LIKE '%int%'
  AND c.column_key = 'PRI'
ORDER BY c.table_schema, c.table_name, c.column_name`
	SQLMaxPrimaryKey = "SELECT COALESCE(MAX(%s), 0) AS max FROM `%s`.`%s`"
)

type Plugin struct {
	Counter int64
}

var plugin = new(Plugin)

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
			"overflow": cnf.Inputs.MySQL[host].Overflow.Enable,
			"interval": cnf.Inputs.MySQL[host].Overflow.Interval,
			"counter":  p.Counter,
		})

		if !cnf.Inputs.MySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.MySQL[host].Overflow.Enable {
			continue
		}

		if !p.isTimeToCollect(cnf.Inputs.MySQL[host].Overflow.Interval) {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
		})

		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		m.Connect()
		m.FetchAll(SQLOverflow, func(table map[string]string) {
			log.DebugWithFields(name, log.Fields{
				"hostname":  cnf.Inputs.MySQL[host].Hostname,
				"database":  table["TABLE_SCHEMA"],
				"table":     table["TABLE_NAME"],
				"column":    table["COLUMN_NAME"],
				"data_type": table["data_type"],
			})

			m.FetchAll(
				fmt.Sprintf(
					SQLMaxPrimaryKey,
					table["COLUMN_NAME"],
					table["TABLE_SCHEMA"],
					table["TABLE_NAME"],
				), func(primaryKey map[string]string) {
					if value, ok := mysql.ParseNumberValue(primaryKey["max"]); ok {
						c := Column{}
						c.dataType = table["data_type"]
						c.current = value
						c.Unsigned()
						c.Maximum()
						c.Percentage()

						mtc.Add(metrics.Metric{
							Key: "mysql_overflow",
							Tags: []metrics.Tag{
								{Name: "hostname", Value: cnf.Inputs.MySQL[host].Hostname},
								{Name: "schema", Value: table["TABLE_SCHEMA"]},
								{Name: "table", Value: table["TABLE_NAME"]},
								{Name: "data_type", Value: c.dataType},
								{Name: "unsigned", Value: strconv.FormatBool(c.unsigned)}},
							Values: []metrics.Value{
								{Key: "percentage", Value: c.percent},
							},
						})

						log.DebugWithFields(name, log.Fields{
							"hostname":  cnf.Inputs.MySQL[host].Hostname,
							"schema":    table["TABLE_SCHEMA"],
							"table":     table["TABLE_NAME"],
							"column":    table["COLUMN_NAME"],
							"data_type": c.dataType,
							"unsigned":  c.unsigned,
							"value":     value,
							"current":   c.current,
							"maximum":   c.maximum,
							"percent":   c.percent,
						})
					}
				})
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
	c.percent = math.Percentage(c.current, c.maximum)
}

func (p *Plugin) isTimeToCollect(i int) bool {
	if p.Counter == 0 || int(time.Since(time.Unix(p.Counter, 0)).Seconds()) >= i {
		(*p).Counter = int64(time.Now().Unix())

		return true
	}

	return false
}

func init() {
	inputs.Add("InputMySQLOverflow", func() inputs.Input { return plugin })
}
