package slack

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/slack"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
	"github.com/swapbyt3s/zenit/plugins/outputs"
)

type OutputSlack struct {}

func (l *OutputSlack) Collect() {
	if config.File.Slack.Enable {
		for _, key := range checks.Load().Keys() {
			var check  = checks.Load().Exist(key)
			var color  = ""
			var status = ""

			if check.Notify() && check.Status >= checks.Notified {
				if check.Between(check.Value) == checks.Warning {
					color = "warning"
					status = "Warning"
				} else if check.Between(check.Value) == checks.Critical {
					color = "danger"
					status = "Critical"
				} else if check.Status == checks.Recovered {
					color = "good"
					status = "Recovered"
				}

				msg := &slack.Message{
					Channel: config.File.Slack.Channel,
				}

				msg.Add(&slack.Attachment{
					Color: color,
					Text: fmt.Sprintf(
						"*[%s]* %s\n*Hostname:* %s (%s)\n%s",
						status,
						check.Name,
						config.File.General.Hostname,
						config.IPAddress,
						check.Message,
					),
				})

				log.Debug(fmt.Sprintf("Slack - Send event notification for %s with status %s and value %d.", check.Name, status, check.Value))

				msg.Send()
			}
		}
	}
}

func init() {
	outputs.Add("OutputSlack", func() outputs.Output { return &OutputSlack{} })
}
