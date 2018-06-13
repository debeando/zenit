package proxysql

import (
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/lib"
)

type Query struct {
  schema string
  digest string
  count  uint
  sum    uint
  min    uint
  max    uint
}

const QUERY_SQL = `
SELECT schemaname,
       digest_text,
       count_star,
       sum_time,
       min_time,
       max_time
FROM stats.stats_mysql_query_digest;
`

func GetQueries() {
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
    var query Query

    rows.Scan(
      &query.schema,
      &query.digest,
      &query.count,
      &query.sum,
      &query.min,
      &query.max)

    Parser(query)
  }
}
