package proxysql

import (
  "sort"
)

type Stat struct {
  schema  string
  table   string
  command string
  count   uint
  sum     uint
  min     uint
  max     uint
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

func (stats *Stats) Contains(s Stat) bool {
  for i := range(stats.Items) {
    if (stats.Items[i].schema == s.schema && stats.Items[i].table == s.table && stats.Items[i].command == s.command) {
      return true
    }
  }
  return false
}

func (stats *Stats) Increment(s Stat) {
  for i := range(stats.Items) {
    if (stats.Items[i].schema == s.schema && stats.Items[i].table == s.table && stats.Items[i].command == s.command) {
      stats.Items[i].count =+ s.count
      stats.Items[i].sum   =+ s.sum
      stats.Items[i].min   =+ s.min
      stats.Items[i].max   =+ s.max
      break
    }
  }
}

func (stats *Stats) Sort() {
  sort.Sort(BySchemaAndTable(stats.Items))
}
