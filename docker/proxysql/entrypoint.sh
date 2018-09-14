#!/bin/bash
# encoding: UTF-8
set -e

proxysql --initial  -f -c /etc/proxysql.cnf
