package outputs

import (
	"sync"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/outputs/prometheus"
	"github.com/swapbyt3s/zenit/plugins/outputs/slack"
)

func Plugins(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
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
