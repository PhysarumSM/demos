FROM ubuntu:16.04

WORKDIR /app

COPY calc-server .
COPY proxy .
COPY perf.conf .

ENV PROXY_IP=127.0.0.1
ENV PROXY_PORT=8080
ENV SERVICE_PORT=4201

CMD ./proxy $SERVICE_PORT calc-server:1.0 $PROXY_IP:$PROXY_PORT > proxy.log 2>&1 & ./calc-server $PROXY_PORT > server.log 2>&1
