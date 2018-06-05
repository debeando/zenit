// Package collect all stats_mysql_connection_pool stats from ProxySQL.
package dbstats

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/lib"
)

type Connection struct {
  hostgroup       int
  srv_host        string
  srv_port        int
  status          string
  connused        int
  connfree        int
  connok          int
  connerr         int
  queries         int
  bytes_data_sent int
  bytes_data_recv int
  latency_us      int
}

func Connections() {
  fmt.Printf("--> proxysql.stats.connections.\n")

  proxysql_conn, _ := lib.Connect(config.DSN_SRC_PROXYSQL)
  defer proxysql_conn.Close()

  mysql_conn, _ := lib.Connect(config.DSN_DST_MYSQL)
  defer mysql_conn.Close()

  connections := getConnections(proxysql_conn)
  putConnections(mysql_conn, connections)
}

func getConnections(conn *sql.DB) []Connection {
  sql := `
SELECT hostgroup,
       srv_host,
       srv_port,
       status,
       connused,
       connfree,
       connok,
       connerr,
       queries,
       bytes_data_sent,
       bytes_data_recv,
       latency_us
FROM stats.stats_mysql_connection_pool;
`
  connections := []Connection{}
  rows, _     := conn.Query(sql)
  defer rows.Close()

  for rows.Next() {
    var connection Connection

    rows.Scan(
      &connection.hostgroup,
      &connection.srv_host,
      &connection.srv_port,
      &connection.status,
      &connection.connused,
      &connection.connfree,
      &connection.connok,
      &connection.connerr,
      &connection.queries,
      &connection.bytes_data_sent,
      &connection.bytes_data_recv,
      &connection.latency_us)

    connections = append(connections, connection)
  }

  return connections
}

func putConnections(conn *sql.DB, connections []Connection) {
  sql := `
INSERT IGNORE INTO connections (
  host_id,
  hostgroup_id,
  hostname,
  port,
  status,
  conn_used,
  conn_free,
  conn_ok,
  conn_err,
  queries,
  bytes_data_sent,
  bytes_data_recv,
  latency_us
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
`

  for _, c := range  connections {
    _, err := conn.Exec(sql,
      config.HOST_TOKEN,
      c.hostgroup,
      c.srv_host,
      c.srv_port,
      c.status,
      c.connused,
      c.connfree,
      c.connok,
      c.connerr,
      c.queries,
      c.bytes_data_sent,
      c.bytes_data_recv,
      c.latency_us)

    if err != nil {
      panic(err)
    }
  }
}
