package commands

import (
	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLCommandCounters = "SELECT * FROM stats_mysql_commands_counters;"

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	for host := range cnf.Inputs.ProxySQL {
		log.DebugWithFields(name, log.Fields{
			"hostname": cnf.Inputs.ProxySQL[host].Hostname,
			"enable":   cnf.Inputs.ProxySQL[host].Enable,
			"commands": cnf.Inputs.ProxySQL[host].Commands,
		})

		if !cnf.Inputs.ProxySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.ProxySQL[host].Commands {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.ProxySQL[host].Hostname,
		})

		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		m.Connect()
		m.FetchAll(SQLCommandCounters, func(row map[string]string) {
			log.DebugWithFields(name, log.Fields{
				"total_time_us": cast.StringToInt64(row["Total_Time_us"]),
				"total_cnt":     cast.StringToInt64(row["Total_cnt"]),
				"cnt_100us":     cast.StringToInt64(row["cnt_100us"]),
				"cnt_500us":     cast.StringToInt64(row["cnt_500us"]),
				"cnt_1ms":       cast.StringToInt64(row["cnt_1ms"]),
				"cnt_5ms":       cast.StringToInt64(row["cnt_5ms"]),
				"cnt_10ms":      cast.StringToInt64(row["cnt_10ms"]),
				"cnt_50ms":      cast.StringToInt64(row["cnt_50ms"]),
				"cnt_100ms":     cast.StringToInt64(row["cnt_100ms"]),
				"cnt_500ms":     cast.StringToInt64(row["cnt_500ms"]),
				"cnt_1s":        cast.StringToInt64(row["cnt_1s"]),
				"cnt_5s":        cast.StringToInt64(row["cnt_5s"]),
				"cnt_10s":       cast.StringToInt64(row["cnt_10s"]),
				"cnt_infs":      cast.StringToInt64(row["cnt_infs"]),
				"hostname":      cnf.Inputs.ProxySQL[host].Hostname,
			})

			mtc.Add(metrics.Metric{
				Key: "proxysql_commands",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: cnf.Inputs.ProxySQL[host].Hostname},
					{Name: "name", Value: row["Command"]},
				},
				Values: []metrics.Value{
					{Key: "total_time_us", Value: cast.StringToInt64(row["Total_Time_us"])},
					{Key: "total_cnt", Value: cast.StringToInt64(row["Total_cnt"])},
					{Key: "cnt_100us", Value: cast.StringToInt64(row["cnt_100us"])},
					{Key: "cnt_500us", Value: cast.StringToInt64(row["cnt_500us"])},
					{Key: "cnt_1ms", Value: cast.StringToInt64(row["cnt_1ms"])},
					{Key: "cnt_5ms", Value: cast.StringToInt64(row["cnt_5ms"])},
					{Key: "cnt_10ms", Value: cast.StringToInt64(row["cnt_10ms"])},
					{Key: "cnt_50ms", Value: cast.StringToInt64(row["cnt_50ms"])},
					{Key: "cnt_100ms", Value: cast.StringToInt64(row["cnt_100ms"])},
					{Key: "cnt_500ms", Value: cast.StringToInt64(row["cnt_500ms"])},
					{Key: "cnt_1s", Value: cast.StringToInt64(row["cnt_1s"])},
					{Key: "cnt_5s", Value: cast.StringToInt64(row["cnt_5s"])},
					{Key: "cnt_10s", Value: cast.StringToInt64(row["cnt_10s"])},
					{Key: "cnt_infs", Value: cast.StringToInt64(row["cnt_infs"])},
				},
			})
		})
	}
}

func init() {
	inputs.Add("InputProxySQLCommands", func() inputs.Input { return &Plugin{} })
}
