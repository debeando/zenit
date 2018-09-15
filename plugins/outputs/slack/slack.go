// TODO: Write in file.

package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
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
		for _, key := range alerts.Load().Keys() {
			var check  = alerts.Load().Exist(key)
			var color  = ""
			var status = ""

			// log.Printf("D! - Output:Slack - check=%#v\n", check)

			if check.Notify() {
				switch check.Status {
				case alerts.Warning:
					log.Printf("D! - Slack:Send event notification - Warning.\n")
					check.Status = alerts.Warning
					color = "warning"
					status = "Warning"
				case alerts.Critical:
					log.Printf("D! - Slack:Send event notification - Critical.\n")
					check.Status = alerts.Critical
					color = "danger"
					status = "Critical"
				case alerts.Resolved:
					log.Printf("D! - Slack:Send event notification - Resolved.\n")
					alerts.Load().Delete(key)
					color = "good"
					status = "Resolved"
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
		log.Print(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Print(err)
	}

	return resp.StatusCode
}
