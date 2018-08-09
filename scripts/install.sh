#!/bin/bash

set -e

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit
fi

if ! [[ "${OSTYPE}" == "linux"* ]]; then
  echo "Only works on Linux amd64."
  exit
fi

FILE="zenit-linux_amd64.tar.gz"
TAG=$(curl --silent "https://api.github.com/repos/swapbyt3s/zenit/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -f /usr/local/bin/zenit ]; then
  rm -f /usr/local/bin/zenit
fi

if [ -f /usr/bin/zenit ]; then
  rm -f /usr/bin/zenit
fi

if [ ! -z "${FILE}" ]; then
  wget -qO- https://github.com/swapbyt3s/zenit/releases/download/${TAG}/${FILE} | tar xz -C /usr/local/bin/
fi

if [ -f /usr/local/bin/zenit ]; then
  ln -s /usr/local/bin/zenit /usr/bin/zenit
fi

if [ -f /etc/zenit/zenit.ini ]; then
  mkdir -p /etc/zenit/
  curl -s https://raw.githubusercontent.com/swapbyt3s/zenit/master/zenit.ini > /etc/zenit/zenit.ini
fi
