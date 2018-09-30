#!/bin/bash
# encoding: UTF-8
set -e

while :; do
  cat /root/slow.log >> /var/lib/mysql/slow.log
  sleep 0.1
done
