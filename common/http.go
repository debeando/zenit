package common

import (
  "net/http"
  "strings"
)

func HTTPPost(uri string, data string) bool {
  req, err := http.NewRequest(
    "POST",
    uri,
    strings.NewReader(data),
  )

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return false
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return false
  }

  return true
}
