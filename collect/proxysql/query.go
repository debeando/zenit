package proxysql

import (
  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/lib"
)

type Query struct {
  hostgroup   int
  schemaname  string
  digest_text string
  count_star  int
  sum_time    int
  min_time    int
  max_time    int
}

const QUERY_SQL = `
SELECT hostgroup,
       schemaname,
       digest_text,
       count_star,
       sum_time,
       min_time,
       max_time
FROM queries;
`

func GetQueries() (queries []Query) {
  conn, err := lib.Connect(config.DSN_PROXYSQL)
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
    var query Query

    rows.Scan(
      &query.hostgroup,
      &query.schemaname,
      &query.digest_text,
      &query.count_star,
      &query.sum_time,
      &query.min_time,
      &query.max_time)

    queries = append(queries, query)
  }

  return queries
}
