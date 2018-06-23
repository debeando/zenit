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
(1, 'test', 'root', '0x7721D69250CB40'  , 'SELECT c FROM sbtest3 WHERE id=?', 8122800, 1441091306, 1441101551, 7032352665, 100, 1000),
(2, 'test', 'root', '0x3BC2F7549D058B6F', 'SELECT c FROM sbtest4 WHERE id=?', 8100134, 1441091306, 1441101551, 7002512958, 100, 2000),
(2, 'test', 'root', '0x4BC2F7549D08B6H',  'SELECT c FROM sbtest4 WHERE id=?', 8100134, 1441091306, 1441101551, 7002512958, 100, 2000);

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
