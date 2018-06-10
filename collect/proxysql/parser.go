package proxysql

import (
  "regexp"
  "strings"
  "sort"
)

const REGEX_SQL = `^(?i)(SELECT|INSERT|UPDATE|DELETE)(?:.*FROM|.*INTO)?\W+([a-zA-Z.]+)`
var re *regexp.Regexp

func init() {
  re, _ = regexp.Compile(REGEX_SQL)
}

func Parser(queries []Query) {
  stats := LoadStats()

  for _, q := range queries {
    if len(q.digest_text) > 0 {
      schema := q.schemaname
      table, attribute := Match(q.digest_text)

      if ! stats.Contains(schema, table, attribute) {
        stats.AddItem(Stat{
          schema: schema,
          table: table,
          attribute: attribute,
          count: 1,
        })
      } else {
        stats.Increment(schema, table, attribute)
      }
    }
  }

  sort.Sort(BySchemaAndTable(stats.Items))
}

func Match(query string) (table string, attribute string) {
  groups    := re.FindStringSubmatch(query)
  table      = GetTable(groups)
  attribute  = GetAttribute(groups)

  return table, attribute
}

func IsSet(arr []string, index int) bool {
  return (len(arr) > index)
}

func GetAttribute(values []string) string {
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
