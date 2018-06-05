// Package main all servers from ProxySQL and ingest into MySQL.
package dbmain

import (
  "fmt"
  "strings"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/swapbyt3s/proxysql_crawler/config"
  "github.com/swapbyt3s/proxysql_crawler/lib"
)

type Server struct {
  hostgroup_id        int
  hostname            string
  port                int
  status              string
  weight              int
  compression         int
  max_connections     int
  max_replication_lag int
  use_ssl             int
  max_latency_ms      int
  comment             string
}

func Servers() {
  fmt.Printf("--> proxysql.main.servers.\n")

  proxysql_conn, _ := lib.Connect(config.DSN_SRC_PROXYSQL)
  defer proxysql_conn.Close()

  mysql_conn, _ := lib.Connect(config.DSN_DST_MYSQL)
  defer mysql_conn.Close()

  servers := getServers(proxysql_conn)
  putServers(mysql_conn, servers)
  delServers(mysql_conn, servers)
}

func getServers(conn *sql.DB) []Server {
  sql := `
SELECT hostgroup_id,
       hostname,
       port,
       status,
       weight,
       compression,
       max_connections,
       max_replication_lag,
       use_ssl,
       max_latency_ms,
       comment
FROM main.mysql_servers;
`
  servers := []Server{}
  rows, _ := conn.Query(sql)
  defer rows.Close()

  for rows.Next() {
    var server Server

    rows.Scan(
      &server.hostgroup_id,
      &server.hostname,
      &server.port,
      &server.status,
      &server.weight,
      &server.compression,
      &server.max_connections,
      &server.max_replication_lag,
      &server.use_ssl,
      &server.max_latency_ms,
      &server.comment)

    servers = append(servers, server)
  }

  return servers
}

func putServers(conn *sql.DB, servers []Server) {
  sql := `
INSERT IGNORE INTO servers (
  host_id,
  hostgroup_id,
  hostname,
  port,
  status,
  weight,
  compression,
  max_connections,
  max_replication_lag,
  use_ssl,
  max_latency_ms,
  comment
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
`

  for _, s := range  servers {
    // fmt.Printf("--> hostname: %s\n", s.hostname)
    _, err := conn.Exec(sql,
      config.HOST_TOKEN,
      s.hostgroup_id,
      s.hostname,
      s.port,
      s.status,
      s.weight,
      s.compression,
      s.max_connections,
      s.max_replication_lag,
      s.use_ssl,
      s.max_latency_ms,
      s.comment)

    if err != nil {
      panic(err)
    }
  }
}

func delServers(conn *sql.DB, org []Server) {
  srv := []string{}
  sql := `
DELETE FROM servers
WHERE (host_id, hostgroup_id, hostname, port) NOT IN (%s);
`

  for _, o := range org {
    srv = append(srv, fmt.Sprintf("('%s',%d,'%s',%d)", config.HOST_TOKEN, o.hostgroup_id, o.hostname, o.port))
  }

  notinsrvs    := strings.Join(srv,",")
  sqlStatement := fmt.Sprintf(sql, notinsrvs)

  _, err := conn.Exec(sqlStatement)
  if err != nil {
    panic(err)
  }
}
