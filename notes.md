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
