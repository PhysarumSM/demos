FROM ubuntu:16.04

WORKDIR /app

COPY covid19-db .
COPY proxy .
COPY perf.conf .
COPY digicert_github.crt /usr/local/share/ca-certificates/digicert_github.crt

RUN apt-get update && apt-get install -y ca-certificates
RUN update-ca-certificates

ENV PROXY_PORT=4201

ENV PROXY_IP=127.0.0.1
ENV SERVICE_PORT=8080

CMD ./proxy $PROXY_PORT covid19-db:1.0 $PROXY_IP:$SERVICE_PORT > proxy.log 2>&1 & ./covid19-db $PROXY_IP $SERVICE_PORT > db.log 2>&1
