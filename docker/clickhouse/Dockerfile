###############################################################################
# Dockerfile to build ClickHouse container images
# Based on Ubuntu 18.04
###############################################################################

# Set the base image to Ubuntu 18.04
FROM ubuntu:18.04

# File Author / Maintainer
MAINTAINER Swapbyt3s
LABEL vendor="Swapbyt3s" \
      description="ClickHouse 18.10.3 on Centos 18.04" \
      version="18.10.3"

# Update the repository sources list
RUN apt-get update

###############################################################################
# BEGIN INSTALLATION
###############################################################################
# -----------------------------------------------------------------------------
# Install additional packages
# -----------------------------------------------------------------------------
RUN apt-get install -y apt-transport-https dirmngr

# -----------------------------------------------------------------------------
# Add repository package and add repository key
# -----------------------------------------------------------------------------
RUN mkdir -p /etc/apt/sources.list.d && \
    apt-key adv --keyserver keyserver.ubuntu.com --recv E0C56BD4 && \
    echo "deb http://repo.yandex.ru/clickhouse/deb/stable/ main/" | \
    tee /etc/apt/sources.list.d/clickhouse.list && \
    apt-get update

# -----------------------------------------------------------------------------
# Install ClickHouse
# -----------------------------------------------------------------------------
RUN env DEBIAN_FRONTEND=noninteractive \
    apt-get install --allow-unauthenticated -y \
            clickhouse-server=18.10.3 \
            clickhouse-common-static=18.10.3 \
            libgcc-7-dev && \
    env DEBIAN_FRONTEND=noninteractive \
    apt-get install --allow-unauthenticated -y \
    clickhouse-client=18.10.3 \
    clickhouse-common-static=18.10.3 \
    locales \
    tzdata

# -----------------------------------------------------------------------------
# Copy script
# -----------------------------------------------------------------------------
COPY config.xml /etc/clickhouse-server/config.d/
RUN chown -R clickhouse /etc/clickhouse-server/

# -----------------------------------------------------------------------------
# Clear
# -----------------------------------------------------------------------------
RUN rm -rf /var/lib/apt/lists/* /var/cache/debconf && \
    apt-get clean
############################## INSTALLATION END ###############################

# Configure locale
RUN locale-gen en_US.UTF-8
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8

# Define user
USER clickhouse

# Expose the default ClickHouse ports
EXPOSE 9000 8123 9009

# Set the working directory to /root
WORKDIR /root

# Start service on run container
ENTRYPOINT exec /usr/bin/clickhouse-server --config=/etc/clickhouse-server/config.xml
