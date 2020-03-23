FROM ubuntu:16.04

WORKDIR /app

COPY helloworldserver .
COPY proxy .

ENV PROXY_IP=127.0.0.1
ENV PROXY_PORT=8080
ENV SERVICE_PORT=4201

CMD ./proxy $SERVICE_PORT hello-world-server $PROXY_IP:$PROXY_PORT > proxy.log & ./helloworldserver $PROXY_PORT > server.log