package alerts

import (
	"sync"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"

	_ "github.com/swapbyt3s/zenit/plugins/alerts/mysql/connections"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/mysql/lagging"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/mysql/readonly"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/mysql/replication"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/os/cpu"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/os/disk"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/os/mem"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/proxysql/errors"
	_ "github.com/swapbyt3s/zenit/plugins/alerts/proxysql/status"
)

func Alerts(wg *sync.WaitGroup) {
	for {
		for key := range loader.Plugins {
			if creator, ok := loader.Plugins[key]; ok {
				c := creator()
				c.Collect()
			}
		}

		time.Sleep(config.File.General.Interval * time.Second)
	}
}
