#!/bin/bash

echo "==> MySQL-Stop Replication: critical"
docker exec -d -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "STOP SLAVE;"
sleep 20

echo "==> MySQL-Start Replication: normal"
docker exec -d -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "START SLAVE;"
sleep 20

echo "==> MySQL-Stop IO Thread Replication: critical"
docker exec -d -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "STOP SLAVE IO_THREAD;"
sleep 20

echo "==> MySQL-Start IO Thread Replication: normal"
docker exec -d -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "START SLAVE IO_THREAD;"
sleep 20

echo "==> MySQL-Stop SQL Thread Replication: critical"
docker exec -d -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "STOP SLAVE SQL_THREAD;"
sleep 20

echo "==> MySQL-Start SQL Thread Replication: normal"
docker exec -d -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "START SLAVE SQL_THREAD;"
sleep 20

echo "==> MySQL-Error Replication: critical"
docker exec -d -i -t -u root zenit_percona_server_primary /usr/bin/mysql -e "CREATE DATABASE zenit;"
docker exec -d -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "DROP DATABASE zenit;"
docker exec -d -i -t -u root zenit_percona_server_primary /usr/bin/mysql -e "DROP DATABASE zenit;"
sleep 20

echo "==> MySQL-Replication: normal"
docker exec -d -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "SET GLOBAL SQL_SLAVE_SKIP_COUNTER = 1; START SLAVE;"
sleep 20

echo "==> MySQL-MaxConnections: warning"
docker exec -d -i -t -u root zenit_percona_server_secondary /bin/sh -c "/usr/bin/mysql -e 'SET GLOBAL max_connections = 10;'"
docker exec -d -i -t -u root zenit_percona_server_secondary /bin/sh -c "/usr/bin/mysql -e 'SELECT SLEEP(60);'"
docker exec -d -i -t -u root zenit_percona_server_secondary /bin/sh -c "/usr/bin/mysql -e 'SELECT SLEEP(60);'"
docker exec -d -i -t -u root zenit_percona_server_secondary /bin/sh -c "/usr/bin/mysql -e 'SELECT SLEEP(60);'"
docker exec -d -i -t -u root zenit_percona_server_secondary /bin/sh -c "/usr/bin/mysql -e 'SELECT SLEEP(60);'"
sleep 20

echo "==> MySQL-MaxConnections: critical"
docker exec -d -i -t -u root zenit_percona_server_secondary /bin/sh -c "/usr/bin/mysql -e 'SELECT SLEEP(60);'"
docker exec -d -i -t -u root zenit_percona_server_secondary /bin/sh -c "/usr/bin/mysql -e 'SELECT SLEEP(60);'"
sleep 20

echo "==> OS-Disk: warning"
docker exec -d -i -t -u root zenit_percona_server_secondary /bin/sh -c "fallocate -l 40G /root/demo"
sleep 20

echo "==> OS-Disk: critical"
docker exec -d -i -t -u root zenit_percona_server_secondary /bin/sh -c "fallocate -l 50G /root/demo"
sleep 20

echo "==> OS-Disk: normal"
docker exec -d -i -t -u root zenit_percona_server_secondary /bin/sh -c "rm /root/demo"
sleep 20
