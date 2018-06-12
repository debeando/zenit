# Zenit

[Zenit](https://en.wikipedia.org/wiki/Zenit_(satellite)) Project is a crawler stats for ProxySQL and MySQL. Zenit is a
russian was spy satellite.

Zenit tool collect stats data and send to ...
- ProxySQL

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
