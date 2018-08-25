FROM ubuntu:18.04

RUN apt-get update && \
    mkdir -p /etc/zenit/ && \
    mkdir -p /var/lib/mysql/ && \
    apt-get clean

COPY docker/sandbox/loop_log_audit.sh /root/loop_log_audit.sh
COPY docker/sandbox/loop_log_slow.sh /root/loop_log_slow.sh
COPY docker/sandbox/entrypoint.sh /root/entrypoint.sh
COPY zenit.ini /etc/zenit/zenit.ini

RUN chmod +x /root/*.sh
RUN touch /var/log/zenit.log
RUN touch /var/lib/mysql/slow.log
RUN touch /var/lib/mysql/audit.log

ENTRYPOINT ["/root/entrypoint.sh"]
