###############################################################################
# Dockerfile to build Sandbox container images
# Based on Ubuntu 18.04
###############################################################################

# Set the base image to Ubuntu 18.04
FROM ubuntu:18.04

# File Author / Maintainer
MAINTAINER Swapbyt3s
LABEL vendor="Swapbyt3s" \
      description="Sandbox on Ubuntu 18.04"

# Update the repository sources list
RUN apt-get update

###############################################################################
# BEGIN INSTALLATION
###############################################################################
# -----------------------------------------------------------------------------
# Install additional packages
# -----------------------------------------------------------------------------
RUN apt-get install -y vim curl wget

# -----------------------------------------------------------------------------
# Copy script
# -----------------------------------------------------------------------------
RUN mkdir -p /etc/zenit/ && \
    mkdir -p /var/lib/mysql/

COPY assets/tests/slow.log /root/slow.log
COPY assets/tests/audit.log /root/audit.log
COPY docker/sandbox/loop_log_audit.sh /root/loop_log_audit.sh
COPY docker/sandbox/loop_log_slow.sh /root/loop_log_slow.sh
COPY docker/sandbox/entrypoint.sh /root/entrypoint.sh
COPY zenit.yaml /etc/zenit/zenit.yaml

RUN chmod +x /root/*.sh
RUN touch /var/lib/mysql/slow.log
RUN touch /var/lib/mysql/audit.log
RUN sed -i s/127.0.0.1:8123/172.20.1.2:8123/ /etc/zenit/zenit.yaml
RUN sed -i s/127.0.0.1:3306/172.20.1.3:3306/ /etc/zenit/zenit.yaml
RUN sed -i s/root@tcp/admin:admin@tcp/ /etc/zenit/zenit.yaml
RUN sed -i s/radminuser:radminpass@tcp\(127.0.0.1:6032\)/proxysql:admin@tcp\(172.20.1.4:6032\)/ /etc/zenit/zenit.yaml

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
