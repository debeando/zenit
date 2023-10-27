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
	if err := cnf.Load(); err != nil {
		return
	}

	for {
		mtc := metrics.Load()
		mtc.Reset()

		for name, plugin := range inputs.Inputs {
			plugin().Collect(name, cnf, mtc)
		}

		for name, plugin := range outputs.Outputs {
			plugin().Deliver(name, cnf, mtc)
		}

		interval := time.Duration(cnf.General.Interval)
		log.DebugWithFields("Wait until next collect metrics", log.Fields{"interval": interval * time.Second})
		time.Sleep(interval * time.Second)
	}
}
