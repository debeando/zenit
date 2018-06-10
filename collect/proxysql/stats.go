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
  sum       int
  min       int
  max       int
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

func (stats *Stats) Contains(s Stat) int {
  for i := range(stats.Items) {
    if (stats.Items[i].schema == s.schema && stats.Items[i].table == s.table && stats.Items[i].attribute == s.attribute) {
      return i
    }
  }
  return 0
}

func (stats *Stats) Increment(s Stat) {
  i := stats.Contains(s)

  stats.Items[i].count =+ s.count
  stats.Items[i].sum =+ s.sum
  stats.Items[i].min =+ s.min
  stats.Items[i].max =+ s.max
}

// PrometheusFormat: Transform data to Prometheus format.
func (stats *Stats) ToArray() (a []string) {
  stats.Sort()
  for i := range(stats.Items) {
    path := fmt.Sprintf("%s.%s.%s", stats.Items[i].schema, stats.Items[i].table, stats.Items[i].attribute)
    a = append(a, fmt.Sprintf("%s.count %d", path, stats.Items[i].count))
    a = append(a, fmt.Sprintf("%s.time %d", path, stats.Items[i].sum))
    a = append(a, fmt.Sprintf("%s.min %d", path, stats.Items[i].min))
    a = append(a, fmt.Sprintf("%s.max %d", path, stats.Items[i].max))
    a = append(a, fmt.Sprintf("%s.avg %d", path, (stats.Items[i].sum / stats.Items[i].count)))
  }
  return a
}

func (stats *Stats) Sort() {
  sort.Sort(BySchemaAndTable(stats.Items))
}
