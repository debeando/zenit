// Package to create schemas on MySQL destionation.
package setup

import (
  "fmt"
  "github.com/swapbyt3s/proxysql_crawler/config"
  "github.com/swapbyt3s/proxysql_crawler/lib"
)

var dropTableStatements = []string {
  "DROP TABLE IF EXISTS hosts;",
  "DROP TABLE IF EXISTS servers;",
  "DROP TABLE IF EXISTS connections;",
  "DROP TABLE IF EXISTS queries;",
}

var createTableStatements = []string {
  `
CREATE TABLE IF NOT EXISTS hosts (
  id CHAR(33) NOT NULL,
  dsn VARCHAR(255) NOT NULL,
  role ENUM('MySQL', 'ProxySQL') NOT NULL DEFAULT 'MySQL',
  UNIQUE (dsn),
  PRIMARY KEY (id)
);`,
  `
CREATE TABLE IF NOT EXISTS servers (
  id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  host_id CHAR(33) NOT NULL,
  hostgroup_id INT NOT NULL DEFAULT 0,
  hostname VARCHAR(255) NOT NULL,
  port INT NOT NULL DEFAULT 3306,
  status ENUM('ONLINE','SHUNNED','OFFLINE_SOFT', 'OFFLINE_HARD') NOT NULL DEFAULT 'ONLINE',
  weight INT NOT NULL DEFAULT 1,
  compression INT NOT NULL DEFAULT 0,
  max_connections INT NOT NULL DEFAULT 1000,
  max_replication_lag INT NOT NULL DEFAULT 0,
  use_ssl INT NOT NULL DEFAULT 0,
  max_latency_ms INT NOT NULL DEFAULT 0,
  comment VARCHAR(255) NOT NULL DEFAULT '',
  UNIQUE (hostgroup_id, hostname, port),
  PRIMARY KEY (id)
);`,
  `
CREATE TABLE IF NOT EXISTS connections (
  id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  host_id CHAR(33) NOT NULL,
  hostgroup_id INT NOT NULL DEFAULT 0,
  hostname VARCHAR(255) NOT NULL,
  port INT NOT NULL DEFAULT 3306,
  status ENUM('ONLINE', 'OFFLINE', 'OFFLINE_SOFT', 'OFFLINE_HARD', 'SHUNNED') NOT NULL,
  conn_used INT NOT NULL DEFAULT 0,
  conn_free INT NOT NULL DEFAULT 0,
  conn_ok INT NOT NULL DEFAULT 0,
  conn_err INT NOT NULL DEFAULT 0,
  queries INT NOT NULL DEFAULT 0,
  bytes_data_sent INT NOT NULL DEFAULT 0,
  bytes_data_recv INT NOT NULL DEFAULT 0,
  latency_us INT NOT NULL DEFAULT 0,
  UNIQUE (hostgroup_id, hostname, port),
  PRIMARY KEY (id)
);`,
`
CREATE TABLE IF NOT EXISTS queries (
  id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  host_id CHAR(33) NOT NULL,
  hostgroup INT NOT NULL DEFAULT 0,
  schemaname VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL,
  digest VARCHAR(64) NOT NULL,
  digest_text TEXT NOT NULL,
  count_star INT NOT NULL DEFAULT 0,
  first_seen INT NOT NULL DEFAULT 0,
  last_seen INT NOT NULL DEFAULT 0,
  sum_time INT NOT NULL DEFAULT 0,
  min_time INT NOT NULL DEFAULT 0,
  max_time INT NOT NULL DEFAULT 0,
  UNIQUE (schemaname, username, digest),
  PRIMARY KEY (id)
);
`,
}

func Run() {
  fmt.Printf("==> Setup...\n")
  mysql_conn, _ := lib.Connect(config.DSN_DST_MYSQL)
  defer mysql_conn.Close()
  fmt.Printf("--> Drop tables...\n")
  lib.ExecCollections(mysql_conn, dropTableStatements)
  fmt.Printf("--> Create tables...\n")
  lib.ExecCollections(mysql_conn, createTableStatements)
}
