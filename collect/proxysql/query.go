package proxysql

import (
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/lib"
)

type Query struct {
  group  string
  schema string
  digest string
  count  uint
  sum    uint
  min    uint
  max    uint
}

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
      &query.group,
      &query.schema,
      &query.digest,
      &query.count,
      &query.sum,
      &query.min,
      &query.max)

    Parser(query)
  }
}
