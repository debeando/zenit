package plugins

import (
	"sync"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"

	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/audit"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/slow"

	_ "github.com/swapbyt3s/zenit/plugins/alerts/mysql/connections"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/mysql/lagging"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/mysql/readonly"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/mysql/replication"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/os/cpu"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/os/disk"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/os/mem"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/proxysql/errors"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/proxysql/status"

	_ "github.com/swapbyt3s/zenit/plugins/inputs/mysql/overflow"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/mysql/slave"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/mysql/status"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/mysql/tables"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/mysql/variables"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/os/cpu"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/os/disk"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/os/mem"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/os/net"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/os/sys"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/percona/deadlock"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/percona/delay"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/percona/kill"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/percona/osc"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/percona/xtrabackup"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/proxysql/commands"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/proxysql/pool"
	_ "github.com/swapbyt3s/zenit/plugins/inputs/proxysql/queries"
)

func Load(wg *sync.WaitGroup) {
	defer wg.Done()

	audit.Start()
	slow.Start()

	for {
		// Flush old metrics:
		metrics.Load().Reset()

		for key := range loader.Plugins {
			if creator, ok := loader.Plugins[key]; ok {
				c := creator()
				c.Collect()
			}
		}

		// Wait loop:
		time.Sleep(config.File.General.Interval * time.Second)
	}
}
