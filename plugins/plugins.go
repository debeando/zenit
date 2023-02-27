package plugins

import (
	"time"

	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/lists/metrics"

	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/outputs"

	_ "github.com/debeando/zenit/plugins/inputs/aws"
	_ "github.com/debeando/zenit/plugins/inputs/mysql/aurora"
	_ "github.com/debeando/zenit/plugins/inputs/mysql/overflow"
	_ "github.com/debeando/zenit/plugins/inputs/mysql/slave"
	_ "github.com/debeando/zenit/plugins/inputs/mysql/status"
	_ "github.com/debeando/zenit/plugins/inputs/mysql/tables"
	_ "github.com/debeando/zenit/plugins/inputs/mysql/variables"
	_ "github.com/debeando/zenit/plugins/inputs/os/cpu"
	_ "github.com/debeando/zenit/plugins/inputs/os/disk"
	_ "github.com/debeando/zenit/plugins/inputs/os/mem"
	_ "github.com/debeando/zenit/plugins/inputs/os/net"
	_ "github.com/debeando/zenit/plugins/inputs/os/sys"
	_ "github.com/debeando/zenit/plugins/inputs/parsers/mysql/audit"
	_ "github.com/debeando/zenit/plugins/inputs/parsers/mysql/slow"
	_ "github.com/debeando/zenit/plugins/inputs/percona/deadlock"
	_ "github.com/debeando/zenit/plugins/inputs/percona/delay"
	_ "github.com/debeando/zenit/plugins/inputs/percona/kill"
	_ "github.com/debeando/zenit/plugins/inputs/percona/osc"
	_ "github.com/debeando/zenit/plugins/inputs/percona/xtrabackup"
	_ "github.com/debeando/zenit/plugins/inputs/proxysql/commands"
	_ "github.com/debeando/zenit/plugins/inputs/proxysql/errors"
	_ "github.com/debeando/zenit/plugins/inputs/proxysql/global"
	_ "github.com/debeando/zenit/plugins/inputs/proxysql/pool"
	_ "github.com/debeando/zenit/plugins/inputs/proxysql/queries"

	_ "github.com/debeando/zenit/plugins/outputs/influxdb"
)

func Load() {
	for {
		// Flush old metrics:
		metrics.Load().Reset()

		for key := range inputs.Inputs {
			if creator, ok := inputs.Inputs[key]; ok {
				c := creator()
				c.Collect()
			}
		}

		for key := range outputs.Outputs {
			if creator, ok := outputs.Outputs[key]; ok {
				c := creator()
				c.Collect()
			}
		}

		// Wait loop:
		time.Sleep(config.File.General.Interval * time.Second)
	}
}
