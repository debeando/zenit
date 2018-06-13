package proxysql

import (
  "regexp"
  "strings"
)

const REGEX_SQL = `^(?i)(SELECT|INSERT|UPDATE|DELETE)(?:.*FROM|.*INTO)?\W+([a-zA-Z0-9._]+)`
var re *regexp.Regexp

func init() {
  re, _ = regexp.Compile(REGEX_SQL)
}

func Parser(q Query) {
  stats := LoadStats()

  if len(q.digest) > 0 {
    table, command := Match(q.digest)

    item := Stat{}
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
