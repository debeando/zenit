package commands

import (
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs/proxysql"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type Command struct {
	Name        string
	TotalTimeUs uint
	TotalCnt    uint
	Cnt100us    uint
	Cnt500us    uint
	Cnt1ms      uint
	Cnt5ms      uint
	Cnt10ms     uint
	Cnt50ms     uint
	Cnt100ms    uint
	Cnt500ms    uint
	Cnt1s       uint
	Cnt5s       uint
	Cnt10s      uint
	CntINFs     uint
}

const SQL = "SELECT * FROM stats_mysql_commands_counters;"

type InputProxySQLCommands struct {}

func (l *InputProxySQLCommands) Collect() {
	if ! config.File.ProxySQL.Inputs.Commands {
		return
	}

	if ! proxysql.Check() {
		return
	}

	conn, err := mysql.Connect(config.File.ProxySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Error("ProxySQL - Impossible to connect: " + err.Error())
	}

	rows, err := conn.Query(SQL)
	defer rows.Close()
	if err != nil {
		log.Error("ProxySQL - Impossible to execute query: " + err.Error())
	}

	for rows.Next() {
		var c Command

		rows.Scan(
			&c.Name,
			&c.TotalTimeUs,
			&c.TotalCnt,
			&c.Cnt100us,
			&c.Cnt500us,
			&c.Cnt1ms,
			&c.Cnt5ms,
			&c.Cnt10ms,
			&c.Cnt50ms,
			&c.Cnt100ms,
			&c.Cnt500ms,
			&c.Cnt1s,
			&c.Cnt5s,
			&c.Cnt10s,
			&c.CntINFs,
		)

		metrics.Load().Add(metrics.Metric{
			Key: "zenit_proxysql_commands",
			Tags: []metrics.Tag{
				{"name", c.Name},
			},
			Values: []metrics.Value{
				{"total_time_us", c.TotalTimeUs},
				{"total_cnt", c.TotalCnt},
				{"cnt_100us", c.Cnt100us},
				{"cnt_500us", c.Cnt500us},
				{"cnt_1ms", c.Cnt1ms},
				{"cnt_5ms", c.Cnt5ms},
				{"cnt_10ms", c.Cnt10ms},
				{"cnt_50ms", c.Cnt50ms},
				{"cnt_100ms", c.Cnt100ms},
				{"cnt_500ms", c.Cnt500ms},
				{"cnt_1s", c.Cnt1s},
				{"cnt_5s", c.Cnt5s},
				{"cnt_10s", c.Cnt10s},
				{"cnt_infs", c.CntINFs},
			},
		})
	}
}

func init() {
	loader.Add("InputProxySQLCommands", func() loader.Plugin { return &InputProxySQLCommands{} })
}
