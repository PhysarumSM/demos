FROM ubuntu:16.04

WORKDIR /app

COPY covid19-webserver .
COPY proxy .
COPY perf.conf .
COPY favicon.ico .

ENV PROXY_PORT=4201

ENV PROXY_IP=127.0.0.1
ENV SERVICE_PORT=8080

CMD ./proxy $PROXY_PORT > proxy.log 2>&1 & ./covid19-webserver $PROXY_PORT $SERVICE_PORT > webserver.log 2>&1
