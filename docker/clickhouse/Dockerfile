###############################################################################
# Dockerfile to build ClickHouse container images
# Based on Ubuntu 18.04
###############################################################################

# Set the base image to Ubuntu 18.04
FROM ubuntu:18.04

# File Author / Maintainer
MAINTAINER Swapbyt3s
LABEL vendor="Swapbyt3s" \
      description="ClickHouse 18.16.1 on Ubuntu 18.04" \
      version="18.16.1"

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
            clickhouse-server=18.16.1 \
            clickhouse-common-static=18.16.1 \
            libgcc-7-dev && \
    env DEBIAN_FRONTEND=noninteractive \
    apt-get install --allow-unauthenticated -y \
    clickhouse-client=18.16.1 \
    clickhouse-common-static=18.16.1 \
    locales \
    tzdata \
    curl

# -----------------------------------------------------------------------------
# Copy script
# -----------------------------------------------------------------------------
COPY docker/clickhouse/entrypoint.sh /root/entrypoint.sh
COPY docker/clickhouse/populate.sh /root/populate.sh
COPY assets/schema/clickhouse/zenit.sql /root/zenit.sql
COPY docker/clickhouse/config.xml /etc/clickhouse-server/config.d/
RUN chown -R clickhouse /etc/clickhouse-server/
RUN chmod a+x /root/entrypoint.sh
RUN chmod a+x /root/populate.sh

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

# Set the working directory to /root
WORKDIR /root

# Expose the default ClickHouse ports
EXPOSE 9000 8123 9009

# Start service on run container
ENTRYPOINT ["/root/entrypoint.sh"]
