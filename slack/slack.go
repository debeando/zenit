package slack

import (
  "bytes"
  "encoding/json"
  "net/http"
  "gitlab.com/swapbyt3s/zenit/config"
)

type Message struct {
  Text        string        `json:"text"`
  Channel     string        `json:"channel,omitempty"`
  UserName    string        `json:"username,omitempty"`
  Attachments []*Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
  Color  string `json:"color,omitempty"`
  Title  string `json:"title,omitempty"`
  Text   string `json:"text,omitempty"`
  Footer string `json:"footer_icon,omitempty"`
}

func (m *Message) AddAttachment(a *Attachment) {
  m.Attachments = append(m.Attachments, a)
}

func Send(msg *Message) int {
  jsonValues, _ := json.Marshal(msg)

  req, err := http.NewRequest(
    "POST",
    "https://hooks.slack.com/services/" + config.SLACK_TOKEN,
    bytes.NewReader(jsonValues),
  )

  if err != nil {
    panic(err)
  }

  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  defer resp.Body.Close()
  if err != nil {
    panic(err)
  }

  return resp.StatusCode
}
