package alerts

import (
	"log"
	"sync"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/connections"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/readonly"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/replication"
)

func Check(wg *sync.WaitGroup) {
	log.Printf("I! - Starting plugin alerts.\n")

	for {
		if config.File.MySQL.Alerts.ReadOnly.Enable {
			readonly.Check()
		}

		if config.File.MySQL.Alerts.Connections.Enable {
			connections.Check()
		}

		if config.File.MySQL.Alerts.Replication.Enable {
			replication.Check()
		}

		time.Sleep(config.File.General.Interval * time.Second)
	}
}
