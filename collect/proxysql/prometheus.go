package proxysql

import (
  "fmt"
)

func PrometheusExport() {
  stats := LoadStats()
  stats.Sort()

  for i := range(stats.Items) {
    path := "proxysql_stats_queries"
    path  = path + fmt.Sprintf("{schema=\"%s\"",  stats.Items[i].schema)
    path  = path + fmt.Sprintf(",table=\"%s\"",   stats.Items[i].table)
    path  = path + fmt.Sprintf(",command=\"%s\"", stats.Items[i].command)
    path  = path + fmt.Sprintf(",group=\"%s\"",   stats.Items[i].group)

    fmt.Printf("%s,calc=\"count\"} %d\n", path, stats.GetCount(i))
    fmt.Printf("%s,calc=\"time\"} %d\n", path, stats.GetSum(i))
    fmt.Printf("%s,calc=\"min\"} %d\n", path, stats.GetMin(i))
    fmt.Printf("%s,calc=\"max\"} %d\n", path, stats.GetMax(i))
    fmt.Printf("%s,calc=\"avg\"} %d\n", path, stats.GetAvg(i))
  }
}
