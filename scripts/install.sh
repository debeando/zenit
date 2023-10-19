#!/bin/bash

# Apache License Version 2.0, January 2004
# https://github.com/debeando/zenit/blob/master/LICENSE.md

set -e

get_arch() {
    # darwin/amd64: Darwin axetroydeMacBook-Air.local 20.5.0 Darwin Kernel Version 20.5.0: Sat May  8 05:10:33 PDT 2021; root:xnu-7195.121.3~9/RELEASE_X86_64 x86_64
    # linux/amd64: Linux test-ubuntu1804 5.4.0-42-generic #46~18.04.1-Ubuntu SMP Fri Jul 10 07:21:24 UTC 2020 x86_64 x86_64 x86_64 GNU/Linux
    a=$(uname -m)
    case ${a} in
        "x86_64" | "amd64" )
            echo "amd64"
        ;;
        "i386" | "i486" | "i586")
            echo "386"
        ;;
        "aarch64" | "arm64" | "arm")
            echo "arm64"
        ;;
        "mips64el")
            echo "mips64el"
        ;;
        "mips64")
            echo "mips64"
        ;;
        "mips")
            echo "mips"
        ;;
        *)
            echo ${NIL}
        ;;
    esac
}

get_os(){
    echo $(uname -s | awk '{print tolower($0)}')
}

OS=$(get_os)
ARCH=$(get_arch)
FILE="_${OS}_${ARCH}.tar.gz"

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit
fi

if ! [[ "${OS}" == "linux" ]]; then
  echo "Only works on Linux."
  exit
fi

if ! type "wget" > /dev/null; then
  echo "The program 'wget' is currently not installed, please install it to continue."
  exit
fi

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

if [ -d "/usr/share/zenit" ]; then
  mkdir /usr/share/zenit
fi

/usr/bin/zenit service --install
/usr/bin/zenit completion bash > /usr/share/zenit/completion
