#!/bin/bash
# encoding: UTF-8
set -e

/root/populate.sh &

/usr/bin/clickhouse-server --config=/etc/clickhouse-server/config.xml
