###############################################################################
# Dockerfile to build Percona Server 5.x container images
# Based on Ubuntu 18.04
###############################################################################

# Set the base image to Ubuntu 18.04
FROM ubuntu:18.04

# File Author / Maintainer
MAINTAINER Swapbyt3s
LABEL vendor="Swapbyt3s" \
      description="Percona Server 5.x on Centos 18.04" \
      version="5.x"

# Update the repository sources list
RUN apt-get update

ENV PERCONA_VERSION 5.5

###############################################################################
# BEGIN INSTALLATION
###############################################################################
RUN apt-get install -y --no-install-recommends \
    apt-transport-https dirmngr apt-utils apt-transport-https ca-certificates \
    libpwquality-tools cracklib-runtime gnupg wget

RUN wget https://repo.percona.com/apt/percona-release_latest.bionic_all.deb >/dev/null 2>&1 && \
    dpkg -i percona-release_latest.bionic_all.deb && \
    apt-get update

RUN env DEBIAN_FRONTEND=noninteractive \
    apt-get install -y percona-server-server-$PERCONA_VERSION && \
    chown -R mysql:0 /var/lib/mysql /var/run/mysqld && \
    chmod 777 /var/run/mysqld

# -----------------------------------------------------------------------------
# Copy script utility
# -----------------------------------------------------------------------------
COPY entrypoint.sh /root/entrypoint.sh
COPY configure.sh /root/configure.sh
RUN chmod a+x /root/entrypoint.sh
RUN chmod a+x /root/configure.sh

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

# Expose the default MySQL port
EXPOSE 3306

# Set the working directory to /root
WORKDIR /root

# Start service on run container
ENTRYPOINT ["/root/entrypoint.sh"]
