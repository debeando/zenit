package example

const File = `---
general:
  hostname: localhost
  interval: 3 # Seconds
  aws_region: ${AWS_REGION}
  aws_access_key_id: ${AWS_ACCESS_KEY_ID}
  aws_secret_access_key: ${AWS_SECRET_ACCESS_KEY}

inputs:
  aws:
    discover:
      enable: false
      # filter: ".*(prd)|(stg).*"
      filter: "."
      username: monitor
      password: monitor
      plugins:
        mysql:
          enable: false
          aurora: false
          innodb: false
          replica: false
          status: false
          variables: false
          overflow:
            enable: false
            interval: 900 # Seconds
          tables:
            enable: false
            interval: 900 # Seconds
    cloudwatch:
      enable: false
  mongodb:
    - hostname: localhost
      dsn: mongodb://localhost:27017
      # dsn: mongodb://user:password@localhost:27017
      enable: false
      serverstatus: false
      collections:
        enable: false
        interval: 900 # Seconds
  mysql:
    - hostname: localhost
      dsn: root@tcp(127.0.0.1:3306)/?timeout=3s
      enable: false
      aurora: false
      innodb: false
      replica: false
      status: false
      variables: false
      overflow:
        enable: false
        interval: 900 # Seconds
      tables:
        enable: false
        interval: 900 # Seconds
  proxysql:
    - hostname: localhost
      dsn: proxysql:admin@tcp(127.0.0.1:6032)/?timeout=3s
      enable: false
      commands: false
      errors: false
      global: false
      pool: false
      queries: false
  os:
    cpu: false
    disk: false
    limits: false
    mem: false
    net: false
  process:
    pt_deadlock_logger: false
    pt_kill: false
    pt_online_schema_change: false
    pt_slave_delay: false
    xtrabackup: false
outputs:
  influxdb:
    enable: false
    url: http://127.0.0.1:8086
    username: zenit
    password: zenit
    database: zenit
`
