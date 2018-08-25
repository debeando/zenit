FROM ubuntu:18.04

ARG repository="deb http://repo.yandex.ru/clickhouse/deb/stable/ main/"
ARG version=18.10.3

RUN apt-get update && \
    apt-get install -y apt-transport-https dirmngr && \
    mkdir -p /etc/apt/sources.list.d && \
    apt-key adv --keyserver keyserver.ubuntu.com --recv E0C56BD4 && \
    echo $repository | tee /etc/apt/sources.list.d/clickhouse.list && \
    apt-get update && \
    env DEBIAN_FRONTEND=noninteractive apt-get install --allow-unauthenticated -y clickhouse-server=$version clickhouse-common-static=$version libgcc-7-dev && \
    env DEBIAN_FRONTEND=noninteractive apt-get install --allow-unauthenticated -y clickhouse-client=$version clickhouse-common-static=$version locales tzdata && \
    rm -rf /var/lib/apt/lists/* /var/cache/debconf && \
    apt-get clean

COPY config.xml /etc/clickhouse-server/config.d/

RUN chown -R clickhouse /etc/clickhouse-server/
RUN locale-gen en_US.UTF-8

ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8

USER clickhouse

EXPOSE 9000 8123 9009

VOLUME /var/lib/clickhouse

ENTRYPOINT exec /usr/bin/clickhouse-server --config=/etc/clickhouse-server/config.xml
