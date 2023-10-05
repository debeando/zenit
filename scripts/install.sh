#!/bin/bash

# Apache License Version 2.0, January 2004
# https://github.com/debeando/zenit/blob/master/LICENSE.md

set -e

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit
fi

if ! [[ "${OSTYPE}" == "linux"* ]]; then
  echo "Only works on Linux."
  exit
fi

if ! type "wget" > /dev/null; then
  echo "The program 'wget' is currently not installed, please install it to continue."
  exit
fi

FILE="zenit-linux_amd64.tar.gz"
TAG=$(wget -qO- "https://api.github.com/repos/debeando/zenit/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -f /etc/systemd/system/zenit.service ]; then
  /usr/bin/zenit service --uninstall
fi

if [ -f /etc/init/zenit.conf ]; then
  /usr/bin/zenit service --uninstall
fi

if [ -f /usr/local/bin/zenit ]; then
  rm -f /usr/local/bin/zenit
fi

if [ -L /usr/bin/zenit ]; then
  rm -f /usr/bin/zenit
fi

if [ -n "${FILE}" ]; then
  wget -qO- "https://github.com/debeando/zenit/releases/download/${TAG}/${FILE}" | tar xz -C /usr/local/bin/
fi

if [ -f /usr/local/bin/zenit ]; then
  ln -s /usr/local/bin/zenit /usr/bin/zenit
fi

if [ ! -f /etc/zenit/zenit.yaml ]; then
  mkdir -p /etc/zenit/
  /usr/bin/zenit --config-example > /etc/zenit/zenit.yaml
fi

if [ -d "/etc/logrotate.d" ]; then
  if [ ! -f /etc/logrotate.d/zenit ]; then
    cat > /etc/logrotate.d/zenit <<EOL
/var/log/zenit.* {
    rotate 7
    daily
    missingok
    notifempty
    copytruncate
    compress
}
EOL
  fi
fi

/usr/bin/zenit service --install
