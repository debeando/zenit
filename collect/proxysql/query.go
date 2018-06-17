package proxysql

import (
  "regexp"
  "sort"
  "strings"
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/lib"
)

type Query struct {
  group   string
  schema  string
  table   string
  command string
  digest  string
  count   uint
  sum     uint
  min     uint
  max     uint
}

type Queries struct {
  Items []Query
}

type BySchemaAndTable []Query

const REGEX_SQL = `^(?i)(SELECT|INSERT|UPDATE|DELETE)(?:.*FROM|.*INTO)?\W+([a-zA-Z0-9._]+)`
const QUERY_SQL = `
SELECT CASE
         WHEN hostgroup IN (SELECT writer_hostgroup FROM main.mysql_replication_hostgroups) THEN 'writer'
         WHEN hostgroup IN (SELECT reader_hostgroup FROM main.mysql_replication_hostgroups) THEN 'reader'
       END AS 'group',
       schemaname,
       digest_text,
       count_star,
       sum_time,
       min_time,
       max_time
FROM stats.stats_mysql_query_digest;
`

var re           *regexp.Regexp
var list_queries *Queries

func init() {
  re, _ = regexp.Compile(REGEX_SQL)
}

func LoadQueries() *Queries {
  if list_queries == nil {
    list_queries = &Queries{}
  }
  return list_queries
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

func (queries *Queries) AddItem(item Query) []Query {
  queries.Items = append(queries.Items, item)
  return queries.Items
}

func (queries *Queries) Count() int {
  return len(queries.Items)
}

func (queries *Queries) Contains(s Query) bool {
  for i := range(queries.Items) {
    if (queries.Items[i].group   == s.group  &&
        queries.Items[i].schema  == s.schema &&
        queries.Items[i].table   == s.table  &&
        queries.Items[i].command == s.command) {
      return true
    }
  }
  return false
}

func (queries *Queries) Increment(s Query) {
  for i := range(queries.Items) {
    if (queries.Items[i].group   == s.group  &&
        queries.Items[i].schema  == s.schema &&
        queries.Items[i].table   == s.table  &&
        queries.Items[i].command == s.command) {
      queries.Items[i].count =+ s.count
      queries.Items[i].sum   =+ s.sum
      queries.Items[i].min   =+ s.min
      queries.Items[i].max   =+ s.max
      break
    }
  }
}

func (queries *Queries) Sort() {
  sort.Sort(BySchemaAndTable(queries.Items))
}

func (queries *Queries) GetSchema(i int) string {
  return queries.Items[i].schema
}

func (queries *Queries) GetTable(i int) string {
  return queries.Items[i].table
}

func (queries *Queries) GetCommand(i int) string {
  return queries.Items[i].command
}

func (queries *Queries) GetGroup(i int) string {
  return queries.Items[i].group
}

func (queries *Queries) GetCount(i int) uint {
  return queries.Items[i].count
}

func (queries *Queries) GetSum(i int) uint {
  return queries.Items[i].sum
}

func (queries *Queries) GetMin(i int) uint {
  return queries.Items[i].min
}

func (queries *Queries) GetMax(i int) uint {
  return queries.Items[i].max
}

func (queries *Queries) GetAvg(i int) uint {
  return queries.Items[i].sum / queries.Items[i].count
}

func GatherQueries() {
  conn, err := lib.MySQLConnect(config.DSN_PROXYSQL)
  defer conn.Close()
  if err != nil {
    panic(err)
  }

  rows, err := conn.Query(QUERY_SQL)
  defer rows.Close()
  if err != nil {
    panic(err)
  }

  for rows.Next() {
    var q Query

    rows.Scan(
      &q.group,
      &q.schema,
      &q.digest,
      &q.count,
      &q.sum,
      &q.min,
      &q.max)

    Parser(q)
  }
}

func Parser(q Query) {
  stats := LoadQueries()

  if len(q.digest) > 0 {
    table, command := Match(q.digest)

    item := Query{}
    item.group   = q.group
    item.schema  = q.schema
    item.table   = table
    item.command = command
    item.count   = q.count
    item.sum     = q.sum
    item.min     = q.min
    item.max     = q.max

    if ! stats.Contains(item) {
      stats.AddItem(item)
    } else {
      stats.Increment(item)
    }
  }
}

func Match(query string) (table string, command string) {
  groups  := re.FindStringSubmatch(query)
  table    = GetTable(groups)
  command  = GetCommand(groups)

  return table, command
}

func IsSet(arr []string, index int) bool {
  return (len(arr) > index)
}

func GetCommand(values []string) string {
  if IsSet(values, 1) {
    return strings.ToLower(values[1])
  }
  return "unknown"
}

func GetTable(values []string) string {
  if IsSet(values, 2) {
    sql_sentence_schema_table := strings.ToLower(values[2])
    sql_sentence_objetcs := strings.Split(sql_sentence_schema_table, ".")
    return sql_sentence_objetcs[len(sql_sentence_objetcs)-1]
  }
  return "unknown"
}
