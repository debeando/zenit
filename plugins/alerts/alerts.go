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

func Check(wg *sync.WaitGroup) {
	for {
		if config.File.OS.Alerts.CPU.Enable {
			cpu.Check()
		}

		if config.File.OS.Alerts.MEM.Enable {
			mem.Check()
		}

		if config.File.OS.Alerts.Disk.Enable {
			disk.Check()
		}

		if config.File.MySQL.Alerts.ReadOnly.Enable {
			readonly.Check()
		}

		if config.File.MySQL.Alerts.Connections.Enable {
			connections.Check()
		}

		if config.File.MySQL.Alerts.Replication.Enable {
			replication.Check()
			lagging.Check()
		}

		if config.File.ProxySQL.Alerts.Status.Enable {
			status.Check()
		}

		if config.File.ProxySQL.Alerts.Errors.Enable {
			errors.Check()
		}

		time.Sleep(config.File.General.Interval * time.Second)
	}
}
