#!/bin/bash

# echo "==> Stop replication!"
# docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "STOP SLAVE;"
# echo "--> Waiting 6 seconds..."
# sleep 6
# echo -n "--> You received the alert? (y/n)? "
# read answer
# if [ "$answer" != "${answer#[Nn]}" ] ;then
#   exit
# fi

# echo "==> Start replication!"
# docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "START SLAVE;"
# echo "-->Waiting 6 seconds..."
# sleep 6
# echo -n "-->You received the alert? (y/n)? "
# read answer
# if [ "$answer" != "${answer#[Nn]}" ] ;then
#   exit
# fi

# echo "==> Stop IO Thread!"
# docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "STOP IO_THREAD;"
# echo "--> Waiting 6 seconds..."
# sleep 6
# echo -n "--> You received the alert? (y/n)? "
# read answer
# if [ "$answer" != "${answer#[Nn]}" ] ;then
#   exit
# fi

# echo "==> Start IO Thread!"
# docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "START IO_THREAD;"
# echo "--> Waiting 6 seconds..."
# sleep 6
# echo -n "--> You received the alert? (y/n)? "
# read answer
# if [ "$answer" != "${answer#[Nn]}" ] ;then
#   exit
# fi

# echo "==> Stop SQL Thread!"
# docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "STOP SQL_THREAD;"
# echo "--> Waiting 6 seconds..."
# sleep 6
# echo -n "--> You received the alert? (y/n)? "
# read answer
# if [ "$answer" != "${answer#[Nn]}" ] ;then
#   exit
# fi

# echo "==> Start SQL Thread!"
# docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "START SQL_THREAD;"
# echo "--> Waiting 6 seconds..."
# sleep 6
# echo -n "--> You received the alert? (y/n)? "
# read answer
# if [ "$answer" != "${answer#[Nn]}" ] ;then
#   exit
# fi

#echo "==> Build replication error!"
#docker exec -i -t -u root zenit_percona_server_primary /usr/bin/mysql -e "CREATE DATABASE zenit;"
#docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "DROP DATABASE zenit;"
#docker exec -i -t -u root zenit_percona_server_primary /usr/bin/mysql -e "DROP DATABASE zenit;"
#echo "--> Waiting 6 seconds..."
#sleep 6
#echo -n "--> You received the alert? (y/n)? "
#read answer
#if [ "$answer" != "${answer#[Nn]}" ] ;then
#  exit
#fi

#echo "==> Skiping replication error!"
#docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "SET GLOBAL SQL_SLAVE_SKIP_COUNTER = 1; START SLAVE;"
#sleep 6
#echo -n "--> You received the alert? (y/n)? "
#read answer
#if [ "$answer" != "${answer#[Nn]}" ] ;then
#  exit
#fi


echo "==> Increment connections to warning!"
docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "SET GLOBAL max_connections = 10;" &
docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "SELECT SLEEP(60);" &
docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "SELECT SLEEP(60);" &
docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "SELECT SLEEP(60);" &
docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "SELECT SLEEP(60);" &
docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "SELECT SLEEP(60);" &
sleep 6

echo "==> Increment connections to critical!"
docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "SELECT SLEEP(60);" &
docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "SELECT SLEEP(60);" &
docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql -e "SELECT SLEEP(60);" &
sleep 6
