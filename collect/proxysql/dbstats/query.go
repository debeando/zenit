// Package collect all stats_mysql_query_digest stats from ProxySQL.
package dbstats

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/lib"
)

type Query struct {
  hostgroup   int
  schemaname  string
  username    string
  digest      string
  digest_text string
  count_star  int
  first_seen  int
  last_seen   int
  sum_time    int
  min_time    int
  max_time    int
}

func Queries() {
  fmt.Printf("--> proxysql.stats.queries.\n")

  proxysql_conn, _ := lib.Connect(config.DSN_SRC_PROXYSQL)
  defer proxysql_conn.Close()

  mysql_conn, _ := lib.Connect(config.DSN_DST_MYSQL)
  defer mysql_conn.Close()

  queries := getQueries(proxysql_conn)
  putQueries(mysql_conn, queries)
}

func getQueries(conn *sql.DB) []Query {
  sql := `
SELECT hostgroup,
       schemaname,
       username,
       digest,
       digest_text,
       count_star,
       first_seen,
       last_seen,
       sum_time,
       min_time,
       max_time
FROM stats.stats_mysql_query_digest;
`
  queries := []Query{}
  rows, _ := conn.Query(sql)
  defer rows.Close()

  for rows.Next() {
    var query Query

    rows.Scan(
      &query.hostgroup,
      &query.schemaname,
      &query.username,
      &query.digest,
      &query.digest_text,
      &query.count_star,
      &query.first_seen,
      &query.last_seen,
      &query.sum_time,
      &query.min_time,
      &query.max_time)

    queries = append(queries, query)
  }

  return queries
}

func putQueries(conn *sql.DB, queries []Query) {
  sql := `
INSERT IGNORE INTO queries (
  host_id,
  hostgroup,
  schemaname,
  username,
  digest,
  digest_text,
  count_star,
  first_seen,
  last_seen,
  sum_time,
  min_time,
  max_time
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
`

  for _, q := range  queries {
    _, err := conn.Exec(sql,
      config.HOST_TOKEN,
      q.hostgroup,
      q.schemaname,
      q.username,
      q.digest,
      q.digest_text,
      q.count_star,
      q.first_seen,
      q.last_seen,
      q.sum_time,
      q.min_time,
      q.max_time)

    if err != nil {
      panic(err)
    }
  }
}
