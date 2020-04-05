FROM ubuntu:16.04

WORKDIR /app

COPY helloworldserver .
COPY proxy .

ENV PROXY_PORT=4201

ENV PROXY_IP=127.0.0.1
ENV SERVICE_PORT=8080

CMD ./proxy $PROXY_PORT hello-world-server $PROXY_IP:$SERVICE_PORT > proxy.log 2>&1 & ./helloworldserver $SERVICE_PORT > server.log 2>&1