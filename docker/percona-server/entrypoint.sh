#!/bin/bash
# encoding: UTF-8

echo '[Entrypoint] Configure MySQL.'
sed -i 's/#SERVER_ID/'"${SERVER_ID}"'/g' /etc/mysql/my.cnf

echo '[Entrypoint] Prepare File System for MySQL.'
DATADIR="/var/lib/mysql"
SOCKETDIR="/var/run/mysqld/"
SOCKET="/var/run/mysqld/mysqld.sock"
CMD=(mysql --protocol=socket -uroot --socket="$SOCKET")

rm -rf "$DATADIR"
mkdir -p "$DATADIR"
mkdir -p "$SOCKETDIR"
chown -R mysql:mysql "$DATADIR"
chown -R mysql:root "$SOCKETDIR"
touch /var/log/mysql/error.log
> /var/log/mysql/error.log

mysql_install_db --user=mysql --basedir=/usr --datadir="$DATADIR"

echo '[Entrypoint] MySQL Initializing database.'
mysqld --skip-networking \
       --socket="$SOCKET" \
       --datadir="$DATADIR" \
       --user=root &
echo '[Entrypoint] MySQL Database initialized.'

for i in {10..0}; do
  if echo 'SELECT 1' | "${CMD[@]}" &> /dev/null; then
    break
  fi
  echo '[Entrypoint] Waiting for server...'
  sleep 1
done

if [ "$i" = 0 ]; then
  echo >&2 '[Entrypoint] Timeout during MySQL init.'
  exit 1
fi

echo "[Entrypoint] Populate TimeZone in MySQL..."
( mysql_tzinfo_to_sql /usr/share/zoneinfo | "${CMD[@]}" --force mysql ) 2> /dev/null
echo "[Entrypoint] End populate TimeZone in MySQL."

echo "[Entrypoint] Create users..."
"${CMD[@]}" <<-EOSQL
SET @@SESSION.SQL_LOG_BIN=0;
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'172.20.1.%' IDENTIFIED BY 'repl';
GRANT ALL ON *.* TO 'admin'@'%' IDENTIFIED BY 'admin' WITH GRANT OPTION;
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP, INDEX ON *.* TO 'sandbox'@'%' IDENTIFIED BY 'sandbox';
GRANT SELECT, REPLICATION CLIENT, SHOW DATABASES, PROCESS ON *.* TO 'monitor'@'%' IDENTIFIED BY 'monitor';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'172.20.1.%' IDENTIFIED BY 'repl';
FLUSH PRIVILEGES;
SET @@SESSION.SQL_LOG_BIN=1;
EOSQL
echo "[Entrypoint] End create users."

if [ "${SERVER_ID}" == "2" ]; then
  for i in {30..0}; do
    if mysqladmin --host=172.20.1.3 --user=repl --password=repl ping &>/dev/null; then
      break
    fi
    echo '[Entrypoint] Waiting for primary server...'
    sleep 1
  done
  if [ "$i" = 0 ]; then
    echo >&2 '[Entrypoint] Timeout for primary server.'
    exit 1
  fi
  echo "[Entrypoint] Connecting to primary node..."
  "${CMD[@]}" <<-EOSQL
SET GLOBAL read_only = 1;
STOP SLAVE;
RESET SLAVE ALL;
CHANGE MASTER TO
  MASTER_HOST="172.20.1.3",
  MASTER_USER="repl",
  MASTER_PASSWORD="repl",
  MASTER_LOG_FILE="mysql_bin.000004",
  MASTER_LOG_POS=107;
START SLAVE;
SHOW SLAVE STATUS\G
EOSQL
  echo "[Entrypoint] End to connect to primary node."
fi

echo "[Entrypoint] MySQL Initialize shutdown..."
mysqladmin shutdown -uroot --socket="$SOCKET"
echo "[Entrypoint] End MySQL initialize shutdown."

echo '[Entrypoint] MySQL init process done. Ready for start up.'
mysqld --user=root
