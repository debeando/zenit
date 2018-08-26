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
# Copy script
# -----------------------------------------------------------------------------
RUN mkdir -p /etc/zenit/ && \
    mkdir -p /var/lib/mysql/

COPY docker/sandbox/loop_log_audit.sh /root/loop_log_audit.sh
COPY docker/sandbox/loop_log_slow.sh /root/loop_log_slow.sh
COPY docker/sandbox/entrypoint.sh /root/entrypoint.sh
COPY zenit.ini /etc/zenit/zenit.ini

RUN chmod +x /root/*.sh
RUN touch /var/log/zenit.log
RUN touch /var/lib/mysql/slow.log
RUN touch /var/lib/mysql/audit.log

# -----------------------------------------------------------------------------
# Clear
# -----------------------------------------------------------------------------
RUN rm -rf /var/lib/apt/lists/* /var/cache/debconf && \
    apt-get clean
# -----------------------------------------------------------------------------
# Clear
# -----------------------------------------------------------------------------
RUN rm -rf /var/lib/apt/lists/* /var/cache/debconf && \
    apt-get clean
############################## INSTALLATION END ###############################

# Set the working directory to /root
WORKDIR /root

# Start service on run container
ENTRYPOINT ["/root/entrypoint.sh"]
