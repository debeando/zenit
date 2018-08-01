// TODO: Write in file.

package prometheus

import (
  "log"
  "fmt"
  "os"
  "strings"

  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/plugins/accumulator"
)

func init() {
  CreateFile()
}

func Run() {
  var a = accumulator.Load()
  var s string

  for _, m := range *a {
    switch m.Values.(type) {
    case int, uint, uint64, float64:
      s = fmt.Sprintf("%s{%s} %s\n", m.Key, getTags(m.Tags), getValue(m.Values))
    case []accumulator.Value:
      for _, i := range m.Values.([]accumulator.Value) {
        s = fmt.Sprintf("%s{%s,type=\"%s\"} %s\n", m.Key, getTags(m.Tags), i.Key, getValue(i.Value))
      }
    }

    if config.General.Debug {
      log.Printf("D! - Prometheus - %s", s)
    }

    //WriteFile(s)
  }
}

func CreateFile() {
  // detect if file exists
  var _, err = os.Stat(config.Prometheus.TextFile)

  // create file if not exists
  if os.IsNotExist(err) {
    var file, err = os.Create(config.Prometheus.TextFile)
    if err != nil { return }
    defer file.Close()
  }
}

func WriteFile(s string) {
  // open file using READ & WRITE permission
  var file, err = os.OpenFile(config.Prometheus.TextFile, os.O_RDWR, 0644)
  if err != nil { return }
  defer file.Close()

  // write some text line-by-line to file
  _, err = file.WriteString(s)
  if err != nil { return }

  // save changes
  err = file.Sync()
  if err != nil { return }
}

func getTags(tags []accumulator.Tag) string {
  s := []string{}
  for t := range(tags) {
    k := tags[t].Name
    v := strings.ToLower(tags[t].Value)
    s = append(s, fmt.Sprintf("%s=\"%s\"", k, v))
  }
  return strings.Join(s,",")
}

func getValue(value interface{}) string {
  switch v := value.(type) {
  case int, uint, uint64:
    return fmt.Sprintf("%d", v)
  case float64:
    return fmt.Sprintf("%.2f", v)
  }

  return "0"
}
