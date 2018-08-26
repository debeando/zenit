#!/bin/bash
# encoding: UTF-8

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
FLUSH PRIVILEGES;
SET @@SESSION.SQL_LOG_BIN=1;
EOSQL
echo "[Entrypoint] End create users."

echo "[Entrypoint] MySQL Initialize shutdown..."
mysqladmin shutdown -uroot --socket="$SOCKET"
echo "[Entrypoint] End MySQL initialize shutdown."

echo '[Entrypoint] MySQL init process done. Ready for start up.'
mysqld --user=root
