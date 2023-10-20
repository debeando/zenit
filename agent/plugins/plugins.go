package plugins

import (
	"time"

	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/outputs"

	_ "zenit/agent/plugins/inputs/aws/cloudwatch/rds"
	_ "zenit/agent/plugins/inputs/aws/discover/rds"
	_ "zenit/agent/plugins/inputs/mongodb/collections"
	_ "zenit/agent/plugins/inputs/mongodb/serverstatus"
	_ "zenit/agent/plugins/inputs/mysql/aurora"
	_ "zenit/agent/plugins/inputs/mysql/innodb"
	_ "zenit/agent/plugins/inputs/mysql/overflow"
	_ "zenit/agent/plugins/inputs/mysql/replica"
	_ "zenit/agent/plugins/inputs/mysql/status"
	_ "zenit/agent/plugins/inputs/mysql/tables"
	_ "zenit/agent/plugins/inputs/mysql/variables"
	_ "zenit/agent/plugins/inputs/os/cpu"
	_ "zenit/agent/plugins/inputs/os/disk"
	_ "zenit/agent/plugins/inputs/os/mem"
	_ "zenit/agent/plugins/inputs/os/net"
	_ "zenit/agent/plugins/inputs/os/sys"
	_ "zenit/agent/plugins/inputs/percona/deadlock"
	_ "zenit/agent/plugins/inputs/percona/delay"
	_ "zenit/agent/plugins/inputs/percona/kill"
	_ "zenit/agent/plugins/inputs/percona/osc"
	_ "zenit/agent/plugins/inputs/percona/xtrabackup"
	_ "zenit/agent/plugins/inputs/proxysql/commands"
	_ "zenit/agent/plugins/inputs/proxysql/errors"
	_ "zenit/agent/plugins/inputs/proxysql/global"
	_ "zenit/agent/plugins/inputs/proxysql/pool"
	_ "zenit/agent/plugins/inputs/proxysql/queries"
	_ "zenit/agent/plugins/outputs/influxdb"

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
