package proxysql

import (
  "github.com/swapbyt3s/zenit/common/log"
  "github.com/swapbyt3s/zenit/common/mysql"
  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/plugins/lists/accumulator"
)

type Pool struct {
  group         string
  serverHost    string
  serverPort    uint
  status        string
  connUsed      uint
  connFree      uint
  connOK        uint
  connERR       uint
  queries       uint
  bytesDataSent uint
  bytesDataRecv uint
  latency       uint
}

const (
  SQLPool = `SELECT CASE
         WHEN hostgroup IN (SELECT writer_hostgroup FROM main.mysql_replication_hostgroups) THEN 'writer'
         WHEN hostgroup IN (SELECT reader_hostgroup FROM main.mysql_replication_hostgroups) THEN 'reader'
       END AS 'group',
       srv_host,
       srv_port,
       status,
       ConnUsed,
       ConnFree,
       ConnOK,
       ConnERR,
       Queries,
       Bytes_data_sent,
       Bytes_data_recv,
       Latency_us
FROM stats.stats_mysql_connection_pool;`
)

func ConnectionPool() {
  conn, err := mysql.Connect(config.File.ProxySQL.DSN)
  defer conn.Close()
  if err != nil {
    log.Error("ProxySQL - Impossible to connect: " + err.Error())
  }

  rows, err := conn.Query(SQLPool)
  defer rows.Close()
  if err != nil {
    log.Error("ProxySQL - Impossible to execute query: " + err.Error())
  }

  for rows.Next() {
    var q Pool

    rows.Scan(
      &q.group,
      &q.serverHost,
      &q.serverPort,
      &q.status,
      &q.connUsed,
      &q.connFree,
      &q.connOK,
      &q.connERR,
      &q.queries,
      &q.bytesDataSent,
      &q.bytesDataRecv,
      &q.latency)

    accumulator.Load().Add(accumulator.Metric{
      Key: "proxysql_connection_pool",
      Tags: []accumulator.Tag{
        {"group", q.group},
        {"host", q.serverHost},
        {"status", q.status},
      },
      Values: []accumulator.Value{
        {"used", q.connUsed},
        {"free", q.connFree},
        {"ok", q.connOK},
        {"errors", q.connERR},
        {"queries", q.queries},
        {"tx", q.bytesDataSent},
        {"rx", q.bytesDataRecv},
        {"latency", q.latency},
      },
    })
  }
}
