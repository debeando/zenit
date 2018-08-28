#!/bin/bash
# encoding: UTF-8
set -e

while :; do
  cat /root/audit.log >> /var/lib/mysql/audit.log
  sleep 0.1
done
