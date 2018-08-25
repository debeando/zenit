#!/bin/bash

while :; do
  cat /root/test_audit.log >> /var/lib/mysql/audit.log
  sleep 0.1
done
