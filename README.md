# Zenit

[Zenit](https://en.wikipedia.org/wiki/Zenit_(satellite)) Project is a crawler stats for ProxySQL and MySQL. Zenit is a
russian was spy satellite.

Zenit tool collect stats data and send to ...
- ProxySQL

The numeric values has represent time has in microseconds.

## ProxySQL

### Configure

Allow remote access:

  mysql -u admin -padmin -h 127.0.0.1 -P 6032
  SET admin-admin_credentials = "admin:admin;radminuser:radminpass";
  LOAD ADMIN VARIABLES TO RUNTIME;

## Prometheus

Integration for Prometheus,


  cp zenit /usr/local/bin/
  * * * * * /usr/local/bin/zenit -collect-proxysql > /usr/local/prometheus/textfile_collector/proxysql_stats_queries.prom

# Todo:
- pgrep -x "mysqld" > /dev/null
  mysql_running_mysqld
- pgrep -f "pt-kill" > /dev/null
  mysql_running_pt_kill
- pgrep -f "pt-deadlock-logger" > /dev/null
  mysql_running_pt_deadlock_logger
- pgrep -f "pt-slave-delay" > /dev/null
  mysql_running_pt_slave_delay
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
