#!/bin/bash
# encoding: UTF-8

SOCKET="/var/run/mysqld/mysqld.sock"
CMD=(mysql --protocol=socket -uroot --socket="$SOCKET")

"${CMD[@]}" <<-EOSQL
SET GLOBAL slow_query_log_always_write_time=0;
SET GLOBAL slow_query_log_file="/var/lib/mysql/slow.log";
SET GLOBAL slow_query_log=ON;
INSTALL PLUGIN audit_log SONAME 'audit_log.so';
EOSQL
