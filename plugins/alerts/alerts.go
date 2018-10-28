package alerts

import (
	"sync"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/connections"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/lagging"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/readonly"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/replication"
	"github.com/swapbyt3s/zenit/plugins/alerts/os/cpu"
	"github.com/swapbyt3s/zenit/plugins/alerts/os/disk"
	"github.com/swapbyt3s/zenit/plugins/alerts/os/mem"
	"github.com/swapbyt3s/zenit/plugins/alerts/proxysql/errors"
	"github.com/swapbyt3s/zenit/plugins/alerts/proxysql/status"
)

func Alerts(wg *sync.WaitGroup) {
	for {
		if config.File.OS.Alerts.CPU.Enable {
			cpu.Collect()
		}

		if config.File.OS.Alerts.MEM.Enable {
			mem.Collect()
		}

		if config.File.OS.Alerts.Disk.Enable {
			disk.Collect()
		}

		if config.File.MySQL.Alerts.ReadOnly.Enable {
			readonly.Collect()
		}

		if config.File.MySQL.Alerts.Connections.Enable {
			connections.Collect()
		}

		if config.File.MySQL.Alerts.Replication.Enable {
			replication.Collect()
			lagging.Collect()
		}

		if config.File.ProxySQL.Alerts.Status.Enable {
			status.Collect()
		}

		if config.File.ProxySQL.Alerts.Errors.Enable {
			errors.Collect()
		}

		time.Sleep(config.File.General.Interval * time.Second)
	}
}
