package common

import (
  "net/http"
  "strings"
  "time"
)

func HTTPPost(uri string, data string) {
  req, err := http.NewRequest(
    "POST",
    uri,
    strings.NewReader(data),
  )

  timeout := time.Duration(1 * time.Second)
  client := &http.Client{
    Timeout: timeout,
  }
  resp, err := client.Do(req)
  defer resp.Body.Close()
  if err != nil {
    panic(err)
  }
}
