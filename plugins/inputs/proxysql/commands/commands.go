package commands

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const query = "SELECT * FROM stats_mysql_commands_counters;"

type InputProxySQLCommands struct {}

func (l *InputProxySQLCommands) Collect() {
	if ! config.File.ProxySQL.Inputs.Commands {
		return
	}

	var a = metrics.Load()
	var p = mysql.GetInstance("proxysql")
	p.Connect(config.File.ProxySQL.DSN)

	rows := p.Query(query)

	for i := range rows {
		a.Add(metrics.Metric{
			Key: "zenit_proxysql_commands",
			Tags: []metrics.Tag{
				{"name", rows[i]["Command"]},
			},
			Values: []metrics.Value{
				{"total_time_us", rows[i]["Total_Time_us"]},
				{"total_cnt", rows[i]["Total_cnt"]},
				{"cnt_100us", rows[i]["cnt_100us"]},
				{"cnt_500us", rows[i]["cnt_500us"]},
				{"cnt_1ms", rows[i]["cnt_1ms"]},
				{"cnt_5ms", rows[i]["cnt_5ms"]},
				{"cnt_10ms", rows[i]["cnt_10ms"]},
				{"cnt_50ms", rows[i]["cnt_50ms"]},
				{"cnt_100ms", rows[i]["cnt_100ms"]},
				{"cnt_500ms", rows[i]["cnt_500ms"]},
				{"cnt_1s", rows[i]["cnt_1s"]},
				{"cnt_5s", rows[i]["cnt_5s"]},
				{"cnt_10s", rows[i]["cnt_10s"]},
				{"cnt_infs", rows[i]["cnt_infs"]},
			},
		})

		log.Debug(fmt.Sprintf("Plugin - InputProxySQLCommands - %#v", rows[i]))
	}
}

func init() {
	loader.Add("InputProxySQLCommands", func() loader.Plugin { return &InputProxySQLCommands{} })
}
