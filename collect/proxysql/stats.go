package proxysql

import (
  "fmt"
  "sort"
)

type Stat struct {
  schema    string
  table     string
  attribute string
  count     int
  time      int
}

type Stats struct {
  Items []Stat
}

type BySchemaAndTable []Stat

var list *Stats

func LoadStats() *Stats {
  if list == nil {
    list = &Stats{}
  }
  return list
}

func (a BySchemaAndTable) Len() int {
  return len(a)
}

func (a BySchemaAndTable) Swap(i, j int) {
  a[i], a[j] = a[j], a[i]
}

func (a BySchemaAndTable) Less(i, j int) bool {
  if a[i].schema < a[j].schema {
    return true
  }
  if a[i].schema > a[j].schema {
    return false
  }
  return a[i].table < a[j].table
}

func (stats *Stats) AddItem(item Stat) []Stat {
  stats.Items = append(stats.Items, item)
  return stats.Items
}

func (stats *Stats) Count() int {
  return len(stats.Items)
}

func (stats *Stats) Contains(schema string, table string, attribute string) bool {
  for i := range(stats.Items) {
    if (stats.Items[i].schema == schema && stats.Items[i].table == table && stats.Items[i].attribute == attribute) {
      return true
    }
  }
  return false
}

func (stats *Stats) Increment(schema string, table string, attribute string, count int, time int) {
  for i := range(stats.Items) {
    if (stats.Items[i].schema == schema && stats.Items[i].table == table && stats.Items[i].attribute == attribute) {
      stats.Items[i].count =+ count
      stats.Items[i].time  =+ time
    }
  }
}

// PrometheusFormat: Transform data to Prometheus format.
func (stats *Stats) ToArray() (a []string) {
  stats.Sort()
  for i := range(stats.Items) {
    path := fmt.Sprintf("%s.%s.%s", stats.Items[i].schema, stats.Items[i].table, stats.Items[i].attribute)
    a = append(a, fmt.Sprintf("%s.count %d", path, stats.Items[i].count))
    a = append(a, fmt.Sprintf("%s.time %d", path, stats.Items[i].time))
  }
  return a
}

func (stats *Stats) Sort() {
  sort.Sort(BySchemaAndTable(stats.Items))
}
