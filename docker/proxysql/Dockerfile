###############################################################################
# Dockerfile to build ProxySQL 1.4 container images
# Based on Ubuntu 18.04
###############################################################################

# Set the base image to Ubuntu 18.04
FROM ubuntu:18.04

# File Author / Maintainer
MAINTAINER Swapbyt3s
LABEL vendor="Swapbyt3s" \
      description="ProxySQL on Ubuntu 18.04" \
      version="1.4"

# Update the repository sources list
RUN apt-get update && \
    apt-get -y upgrade

ENV PERCONA_VERSION 5.5

###############################################################################
# BEGIN INSTALLATION
###############################################################################
# -----------------------------------------------------------------------------
# Install additional packages
# -----------------------------------------------------------------------------
RUN apt-get install -y vim wget htop stress curl lsb

# -----------------------------------------------------------------------------
# Install ProxySQL packages
# -----------------------------------------------------------------------------
RUN apt-get install -y --no-install-recommends \
    apt-transport-https dirmngr apt-utils apt-transport-https ca-certificates \
    libpwquality-tools cracklib-runtime gnupg

RUN wget https://repo.percona.com/apt/percona-release_latest.$(lsb_release -sc)_all.deb >/dev/null 2>&1 && \
    dpkg -i percona-release_latest.$(lsb_release -sc)_all.deb && \
    apt-get update

RUN env DEBIAN_FRONTEND=noninteractive \
    apt-get install -y percona-server-client-$PERCONA_VERSION

RUN env DEBIAN_FRONTEND=noninteractive \
    apt-get install -y proxysql

# -----------------------------------------------------------------------------
# Install Zenit
# -----------------------------------------------------------------------------
RUN mkdir -p /etc/zenit/
COPY zenit.yaml /etc/zenit/zenit.yaml
RUN sed -i 's/localhost/${HOSTNAME}/' /etc/zenit/zenit.yaml
RUN sed -i 's/debug: false/debug: true/' /etc/zenit/zenit.yaml
RUN sed -i 's/127.0.0.1:8123/172.20.1.2:8123/' /etc/zenit/zenit.yaml
RUN sed -i 's/radminuser:radminpass@tcp\(127.0.0.1:6032\)/proxysql:admin@tcp\(172.20.1.4:6032\)/' /etc/zenit/zenit.yaml
RUN sed -i 's/token: xxx\/yyy\/zzz/token: ${ZENIT_SLACK_TOKEN}/' /etc/zenit/zenit.yaml
RUN sed -i 's/interval: 10/interval: 5/' /etc/zenit/zenit.yaml
RUN sed -i 's/duration: 30/duration: 10/' /etc/zenit/zenit.yaml

# -----------------------------------------------------------------------------
# Copy script utility
# -----------------------------------------------------------------------------
COPY docker/proxysql/entrypoint.sh /root/entrypoint.sh
COPY docker/proxysql/proxysql.cnf /etc/proxysql.cnf
RUN chmod a+x /root/entrypoint.sh
RUN chown 0640 /etc/proxysql.cnf
RUN chown root:proxysql /etc/proxysql.cnf

# -----------------------------------------------------------------------------
# Clear
# -----------------------------------------------------------------------------
RUN rm -rf /var/lib/apt/lists/* /var/cache/debconf && \
    apt-get clean
############################## INSTALLATION END ###############################

# Expose the default MySQL port
EXPOSE 3306 6032

# Set the working directory to /root
WORKDIR /root

# Start service on run container
ENTRYPOINT ["/root/entrypoint.sh"]
