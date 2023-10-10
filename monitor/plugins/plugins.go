package plugins

import (
	"time"

	"zenit/config"
	"zenit/monitor/plugins/lists/metrics"

	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/outputs"

	_ "zenit/monitor/plugins/inputs/aws/cloudwatch/rds"
	_ "zenit/monitor/plugins/inputs/aws/discover/rds"
	_ "zenit/monitor/plugins/inputs/mongodb/serverstatus"
	_ "zenit/monitor/plugins/inputs/mysql/aurora"
	_ "zenit/monitor/plugins/inputs/mysql/overflow"
	_ "zenit/monitor/plugins/inputs/mysql/replica"
	_ "zenit/monitor/plugins/inputs/mysql/status"
	_ "zenit/monitor/plugins/inputs/mysql/tables"
	_ "zenit/monitor/plugins/inputs/mysql/variables"
	_ "zenit/monitor/plugins/inputs/os/cpu"
	_ "zenit/monitor/plugins/inputs/os/disk"
	_ "zenit/monitor/plugins/inputs/os/mem"
	_ "zenit/monitor/plugins/inputs/os/net"
	_ "zenit/monitor/plugins/inputs/os/sys"
	_ "zenit/monitor/plugins/inputs/percona/deadlock"
	_ "zenit/monitor/plugins/inputs/percona/delay"
	_ "zenit/monitor/plugins/inputs/percona/kill"
	_ "zenit/monitor/plugins/inputs/percona/osc"
	_ "zenit/monitor/plugins/inputs/percona/xtrabackup"
	_ "zenit/monitor/plugins/inputs/proxysql/commands"
	_ "zenit/monitor/plugins/inputs/proxysql/errors"
	_ "zenit/monitor/plugins/inputs/proxysql/global"
	_ "zenit/monitor/plugins/inputs/proxysql/pool"
	_ "zenit/monitor/plugins/inputs/proxysql/queries"
	_ "zenit/monitor/plugins/outputs/influxdb"

	"github.com/debeando/go-common/log"
)

func Load() {
	cnf := config.GetInstance()
	cnf.Load()

	for {
		// Flush old metrics:
		mtc := metrics.Load()
		mtc.Reset()

		for key := range inputs.Inputs {
			if creator, ok := inputs.Inputs[key]; ok {
				creator().Collect(key, cnf, mtc)
			}
		}

		for key := range outputs.Outputs {
			if creator, ok := outputs.Outputs[key]; ok {
				creator().Deliver(key, cnf, mtc)
			}
		}

		// Wait loop:
		log.DebugWithFields("Wait until next collect metrics", log.Fields{"interval": cnf.General.Interval * time.Second})
		time.Sleep(config.GetInstance().General.Interval * time.Second)
	}
}
