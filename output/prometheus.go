package output

import (
  "fmt"
  "strings"
)

func Prometheus() {
  var a = Load()

  for _, m := range *a {
    switch m.Values.(type) {
    case uint, uint64, float64:
      fmt.Printf("%s{%s} %s\n", m.Key, getTags(m.Tags), getValue(m.Values))
    case []Value:
      for _, i := range m.Values.([]Value) {
        fmt.Printf("%s{%s,type=\"%s\"} %s\n", m.Key, getTags(m.Tags), i.Key, getValue(i.Value))
      }
    }
  }
}

func getTags(tags []Tag) string {
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
  case uint, uint64:
    return fmt.Sprintf("%d", v)
  case float64:
    return fmt.Sprintf("%.2f", v)
  }

  return "0"
}