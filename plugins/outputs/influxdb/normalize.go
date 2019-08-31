package influxdb

import (
  // "fmt"
  // "github.com/swapbyt3s/zenit/common/log"

  // "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

func Normalize(items *metrics.Items) map[string][]map[string]interface{} {
  e := make(map[string][]map[string]interface{})

  for _, i := range *items {
    tags    := make(map[string]string)
    fields := make(map[string]interface{})

    // log.Debug(fmt.Sprintf("Plugin - OutputIndluxDB L4 - %#v", i))

    for _, t := range i.Tags {
      tags[t.Name] = t.Value
    }

    for _, y := range i.Values {
      fields[y.Key] = y.Value
    }

    item := make(map[string]interface{})
    item["tags"] = tags
    item["fields"] = fields

    e[i.Key] = append(e[i.Key], item)
  }

  return e
}
