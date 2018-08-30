#!/bin/bash
# encoding: UTF-8
set -e

echo "Wait for database ready..."
until $(/usr/bin/curl --output /dev/null --silent --fail --data 'SELECT 1' http://127.0.0.1:8123/?database=system); do
  echo "."
  sleep 1
done

echo "Populating database..."
cat /root/zenit.sql | /usr/bin/clickhouse-client --multiquery
echo "Populated database!"
