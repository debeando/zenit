package slack

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
)

type Message struct {
	Channel     string        `json:"channel,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Color string `json:"color,omitempty"`
	Text  string `json:"text,omitempty"`
}

func (m *Message) Add(a *Attachment) {
	m.Attachments = append(m.Attachments, a)
}

func (m *Message) Send() int {
	jsonValues, _ := json.Marshal(m)

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
