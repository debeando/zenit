package proxysql

import (
  "fmt"
)

func PrometheusExport() {
  stats := LoadStats()
  stats.Sort()

  for i := range(stats.Items) {
    path := "proxysql_stats_queries"
    path  = path + fmt.Sprintf("{schema=%s",  stats.Items[i].schema)
    path  = path + fmt.Sprintf(",table=%s",   stats.Items[i].table)
    path  = path + fmt.Sprintf(",command=%s", stats.Items[i].command)

    fmt.Printf("%s,calc=count} %d\n", path, stats.Items[i].count)
    fmt.Printf("%s,calc=time} %d\n", path, stats.Items[i].sum)
    fmt.Printf("%s,calc=min} %d\n", path, stats.Items[i].min)
    fmt.Printf("%s,calc=max} %d\n", path, stats.Items[i].max)
    fmt.Printf("%s,calc=avg} %d\n", path, (stats.Items[i].sum / stats.Items[i].count))
  }
}
