# Install

Manual instalation:

```bash
sudo cp zenit /usr/bin/zenit
sudo mkdir /etc/zenit/
sudo /usr/bin/zenit agent --config-example > /etc/zenit/zenit.yaml
sudo /usr/bin/zenit service --install
sudo systemctl start zenit
sudo systemctl status zenit
```

## Reinstall

```bash
sudo systemctl stop zenit
sudo /usr/bin/zenit service --uninstall
sudo rm /usr/bin/zenit
```

And repeat the install process.

## Test

```bash
/usr/bin/zenit agent --config=/etc/zenit/zenit.yaml -v
```

## Configure

### ProxySQL

```yaml
general:
  hostname: tangerine-prd-proxysql-node072
  interval: 30

inputs:
  proxysql:
    - hostname: tangerine-prd-proxysql-node072
      dsn: proxysql:admin@tcp(127.0.0.1:6032)/?timeout=3s
      enable: true
      commands: true
      errors: true
      global: true
      pool: true
      queries: false
  os:
    cpu: true
    disk: true
    limits: true
    mem: true
    net: true
outputs:
  influxdb:
    enable: true
    url: http://10.50.0.171:8086
    database: zenit
 ```

## Uninstall

```bash
sudo systemctl stop zenit
sudo /usr/bin/zenit service --uninstall
sudo rm /usr/bin/zenit
sudo rm -rf /etc/zenit/
```
