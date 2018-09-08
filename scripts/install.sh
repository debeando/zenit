#!/bin/bash

# Apache License Version 2.0, January 2004
# https://github.com/swapbyt3s/zenit/blob/master/LICENSE.md

set -e

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit
fi

if ! [[ "${OSTYPE}" == "linux"* ]]; then
  echo "Only works on Linux amd64."
  exit
fi

if ! type "wget" > /dev/null; then
  echo "The program 'wget' is currently not installed, please install it to continue."
  exit
fi

FILE="zenit-linux_amd64.tar.gz"
TAG=$(wget -qO- "https://api.github.com/repos/swapbyt3s/zenit/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -f /usr/local/bin/zenit ]; then
  rm -f /usr/local/bin/zenit
fi

if [ -L /usr/bin/zenit ]; then
  rm -f /usr/bin/zenit
fi

if [ ! -z "${FILE}" ]; then
  wget -qO- https://github.com/swapbyt3s/zenit/releases/download/${TAG}/${FILE} | tar xz -C /usr/local/bin/
fi

if [ -f /usr/local/bin/zenit ]; then
  ln -s /usr/local/bin/zenit /usr/bin/zenit
fi

if [ ! -f /etc/zenit/zenit.yaml ]; then
  mkdir -p /etc/zenit/
  wget -qO- "https://raw.githubusercontent.com/swapbyt3s/zenit/master/zenit.yaml" > /etc/zenit/zenit.yaml
fi

exit 0
