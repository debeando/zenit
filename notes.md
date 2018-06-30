CREATE DATABASE main;

CREATE TABLE main.mysql_replication_hostgroups (
writer_hostgroup INT NOT NULL PRIMARY KEY,
reader_hostgroup INT NOT NULL,
UNIQUE (reader_hostgroup)
);

INSERT INTO main.mysql_replication_hostgroups VALUES (1,2);

CREATE TABLE stats.stats_mysql_query_digest (
  hostgroup INT,
  schemaname VARCHAR(64) NOT NULL,
  username VARCHAR(64) NOT NULL,
  digest VARCHAR(24) NOT NULL,
  digest_text VARCHAR(255) NOT NULL,
  count_star BIGINT NOT NULL,
  first_seen BIGINT NOT NULL,
  last_seen BIGINT NOT NULL,
  sum_time BIGINT NOT NULL,
  min_time BIGINT NOT NULL,
  max_time BIGINT NOT NULL,
  PRIMARY KEY(schemaname, username, digest)
);

TRUNCATE stats.stats_mysql_query_digest;
INSERT INTO stats.stats_mysql_query_digest
VALUES
(1, 'test', 'root', '0x7721', 'SELECT c FROM sbtest3 WHERE id=?', 1, 1441091306, 1441101551, 1, 100, 1000),
(2, 'test', 'root', '0x4H20', 'SELECT c FROM sbtest5 WHERE id=?', 2, 1441091306, 1441101551, 1, 100, 2000),
(2, 'test', 'root', '0x4H2A', 'SELECT c FROM sbtest4 WHERE id=?', 2, 1441091306, 1441101551, 1, 100, 2000),
(2, 'test', 'root', '0x4H2B', 'SELECT c FROM sbtest4 WHERE id=?', 2, 1441091306, 1441101551, 1, 100, 2000),
(2, 'test', 'root', '0x4H2C', 'SELECT c FROM sbtest4 WHERE id=?', 1, 1441091306, 1441101551, 1, 100, 2000),
(2, 'test', 'root', '0x4H2D', 'SELECT c FROM sbtest4 WHERE id=?', 3, 1441091306, 1441101551, 1, 100, 2000);

UPDATE global_variables SET variable_value='0.0.0.0:6032' WHERE variable_name='mysql-interfaces';
SET admin-admin_credentials="admin:admin;radminuser:radminpass";

LOAD ADMIN VARIABLES TO RUNTIME;
PROXYSQL RESTART;

CREATE TABLE connection_pool (
  id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  proxysql_id BIGINT(20) UNSIGNED NOT NULL
  hostgroup
  srv_host
  srv_port
  status ENUM('ONLINE', 'OFFLINE', 'OFFLINE_SOFT', 'OFFLINE_HARD', 'SHUNNED') NOT NULL,
  ConnUsed
  ConnFree
  ConnOK
  ConnERR
  Queries
  Bytes_data_sent
  Bytes_data_recv
);

stats_mysql_query_digest
stats_mysql_commands_counters

stats_mysql_connection_pool;

https://github.com/sysown/proxysql/wiki/ProxySQL-Configuration


INSERT INTO mysql_replication_hostgroups (
  writer_hostgroup,
  reader_hostgroup,
  comment
) VALUES (1,2,'cluster');

INSERT INTO mysql_servers (
      hostgroup_id,
      hostname,
      port,
      max_connections,
      max_replication_lag
    ) VALUES (
      1,
      '192.168.1.35',
      3306,
      10,
      60
    );

INSERT INTO mysql_servers (
      hostgroup_id,
      hostname,
      port,
      max_connections,
      max_replication_lag
    ) VALUES (
      1,
      '192.168.1.100',
      3306,
      10,
      60
    );


LOAD MYSQL SERVERS TO RUNTIME;
SAVE MYSQL SERVERS TO DISK;

root@127.0.0.1:proxysql_stats> SELECT * FROM servers WHERE (proxysql_id, hostgroup_id, hostname, port) NOT IN ((1, 1, '192.168.1.35', 3306), (1,1, '192.168.1.100', 3306));


  // -collect-mysql {status,variables,slave status}
  // check is open port from any
  // -output-prometheus
  // -output-influxdb
  // -percona-skip-replication-error
  // -percona-eta-catchup
  // -collect-os {network,swap,disk,iops}
  // -collect-pgbouncer ?
  // -collect-postgresql ?
  // -collect-mongodb ?


  CREATE TABLE IF NOT EXISTS zenit.demo (
    _time DateTime,
    _date Date default toDate(_time),
    value String
  ) ENGINE = MergeTree(_date, (_time), 8192);

  INSERT INTO zenit.demo (value) VALUES ('A');

  SELECT * FROM zenit.demo\G

curl -s -d 'SELECT 1' http://172.17.0.3:8123/?database=zenit

curl -s -X POST -d "INSERT INTO demo (value) VALUES ('A')" http://10.201.17.217:8123/?database=zenit

cat schema/mysql.sql | clickhouse-client --multiline


# Todo:
- @@log_error
  mysql_errors_on_log
# Check if running audit plugin?
# Check if running general log?
# Check if running slow log?
# Check SQL safe:
# - SELECT @@SQL_SAFE_UPDATES;
# - SELECT @@SQL_SELECT_LIMIT;
# - SELECT @@MAX_JOIN_SIZE;
- have log rotation file? for
  > audit log
  > general log
  > error log
  > slow log
- To ClickHouse
  > audit log
  > general log
  > error log
  > slow log
  > binarylogs

https://hub.docker.com/r/yandex/clickhouse-server/

docker run -d --name some-clickhouse-server --ulimit nofile=262144:262144 yandex/clickhouse-server
docker run -it --rm --link some-clickhouse-server:clickhouse-server yandex/clickhouse-client --host clickhouse-server

GOOS=linux go build -ldflags "-s -w" -o zenit main.go && docker cp zenit d1c86f2f36ff:/root && docker exec -i -t d1c86f2f36ff /root/zenit -parser=auditlog-xml -parser-file=/root/test_audit.log
