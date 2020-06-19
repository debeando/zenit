package config

const ExampleFile = `---
general:
  hostname: localhost
  interval: 10 # Seconds
  debug: true
  aws_region: ${AWS_REGION}
  aws_access_key_id: ${AWS_ACCESS_KEY_ID}
  aws_secret_access_key: ${AWS_SECRET_ACCESS_KEY}

parser:
  mysql:
    slowlog:
      enable: true
      log_path: /var/lib/mysql/slow.log
      buffer_size: 100   # Number of events.
      buffer_timeout: 60 # Seconds
    auditlog:
      enable: true
      format: xml-old
      log_path: /var/lib/mysql/audit.log
      buffer_size: 100   # Number of events.
      buffer_timeout: 60 # Seconds

inputs:
  awsdiscover:
    enable: true
    username: monitor
    password: monitor
    filter: prd
    plugins:
      aurora: false
      overflow: true
      slave: true
      status: true
      tables: true
      variables: true
  mysql:
    - hostname: localhost
      dsn: root@tcp(127.0.0.1:3306)/?timeout=3s
      aurora: false
      overflow: true
      slave: true
      status: true
      tables: true
      variables: true
  proxysql:
    - hostname: localhost
      dsn: proxysql:admin@tcp(127.0.0.1:6032)/?timeout=3s
      commands: true
      errors: true
      global: true
      pool: true
      queries: true
  os:
    cpu: true
    disk: true
    limits: true
    mem: true
    net: true
  process:
    pt_deadlock_logger: true
    pt_kill: true
    pt_online_schema_change: true
    pt_slave_delay: true
    xtrabackup: true

outputs:
  clickhouse:
    enable: false
    dsn: http://127.0.0.1:8123/?database=zenit
    # dsn: http://127.0.0.1:8123/?user=admin&password=admin&database=zenit
  influxdb:
    enable: true
    url: http://127.0.0.1:8086
    # username: zenit
    # password: zenit
    database: zenit
`
