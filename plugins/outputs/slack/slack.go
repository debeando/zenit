package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
)

type Message struct {
	Channel     string        `json:"channel,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Color string `json:"color,omitempty"`
	Text  string `json:"text,omitempty"`
}

func (m *Message) AddAttachment(a *Attachment) {
	m.Attachments = append(m.Attachments, a)
}

func Run() {
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

				msg := &Message{
					Channel: config.File.Slack.Channel,
				}

				msg.AddAttachment(&Attachment{
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

				Send(msg)
			}
		}
	}
}

func Send(msg *Message) int {
	jsonValues, _ := json.Marshal(msg)

	req, err := http.NewRequest(
		"POST",
		"https://hooks.slack.com/services/" + config.File.Slack.Token,
		bytes.NewReader(jsonValues),
	)

	if err != nil {
		log.Error(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Error(err.Error())
	}

	return resp.StatusCode
}
