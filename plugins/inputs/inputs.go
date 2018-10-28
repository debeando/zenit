// TODO:
// - Convert this into module/package called "collect" because use for inputs and parsers.
// - If not set any option, ignore and no enter in infinite loop.

package inputs

import (
	"sync"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/audit"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/indexes"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/overflow"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/slave"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/slow"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/status"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/tables"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/variables"
	"github.com/swapbyt3s/zenit/plugins/inputs/os/cpu"
	"github.com/swapbyt3s/zenit/plugins/inputs/os/disk"
	"github.com/swapbyt3s/zenit/plugins/inputs/os/mem"
	"github.com/swapbyt3s/zenit/plugins/inputs/os/net"
	"github.com/swapbyt3s/zenit/plugins/inputs/os/sys"
	"github.com/swapbyt3s/zenit/plugins/inputs/process"
	"github.com/swapbyt3s/zenit/plugins/inputs/proxysql"
	"github.com/swapbyt3s/zenit/plugins/inputs/proxysql/commands"
	"github.com/swapbyt3s/zenit/plugins/inputs/proxysql/pool"
	"github.com/swapbyt3s/zenit/plugins/inputs/proxysql/queries"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/outputs/prometheus"
	"github.com/swapbyt3s/zenit/plugins/outputs/slack"
)

func Plugins(wg *sync.WaitGroup) {
	defer wg.Done()

	audit.Collect()
	slow.Collect()

	for {
		// Flush old metrics:
		metrics.Load().Reset()

		// Collect by OS:
		if config.File.OS.Inputs.CPU {
			cpu.Collect()
		}
		if config.File.OS.Inputs.Disk {
			disk.Collect()
		}
		if config.File.OS.Inputs.Mem {
			mem.Collect()
		}
		if config.File.OS.Inputs.Net {
			net.Collect()
		}
		if config.File.OS.Inputs.Limits {
			sys.Collect()
		}

		// Collect by MySQL:
		if mysql.Check() {
			if config.File.MySQL.Inputs.Indexes {
				indexes.Collect()
			}
			if config.File.MySQL.Inputs.Overflow {
				overflow.Collect()
			}
			if config.File.MySQL.Inputs.Slave {
				slave.Collect()
			}
			if config.File.MySQL.Inputs.Status {
				status.Collect()
			}
			if config.File.MySQL.Inputs.Tables {
				tables.Collect()
			}
			if config.File.MySQL.Inputs.Variables {
				variables.Collect()
			}
		}

		// Collect by ProxySQL:
		if proxysql.Check() {
			if config.File.ProxySQL.Inputs.Commands {
				commands.Collect()
			}
			if config.File.ProxySQL.Inputs.Pool {
				pool.Collect()
			}
			if config.File.ProxySQL.Inputs.Queries {
				queries.Collect()
			}
		}

		// Collect by process:
		if config.File.Process.Inputs.PerconaToolKitKill {
			process.PerconaToolKitKill()
		}
		if config.File.Process.Inputs.PerconaToolKitDeadlockLogger {
			process.PerconaToolKitDeadlockLogger()
		}
		if config.File.Process.Inputs.PerconaToolKitSlaveDelay {
			process.PerconaToolKitSlaveDelay()
		}
		if config.File.Process.Inputs.PerconaToolKitOnlineSchemaChange {
			process.PerconaToolKitOnlineSchemaChange()
		}
		if config.File.Process.Inputs.PerconaXtraBackup {
			process.PerconaXtraBackup()
		}
		if config.File.Prometheus.Enable {
			prometheus.Run()
		}
		if config.File.Slack.Enable {
			slack.Run()
		}

		// Wait loop:
		time.Sleep(config.File.General.Interval * time.Second)
	}
}
