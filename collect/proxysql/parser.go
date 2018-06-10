package proxysql

import (
  "regexp"
  "strings"
)

const REGEX_SQL = `^(?i)(SELECT|INSERT|UPDATE|DELETE)(?:.*FROM|.*INTO)?\W+([a-zA-Z.]+)`
var re *regexp.Regexp

func init() {
  re, _ = regexp.Compile(REGEX_SQL)
}

func Parser(digest Query) {
  stats := LoadStats()

  if len(digest.digest_text) > 0 {
    table, attribute := Match(digest.digest_text)

    if ! stats.Contains(digest.schemaname, table, attribute) {
      stats.AddItem(Stat{
        schema: digest.schemaname,
        table: table,
        attribute: attribute,
        count: digest.count_star,
      })
    } else {
      stats.Increment(digest.schemaname, table, attribute, digest.count_star, digest.sum_time)
    }
  }
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
