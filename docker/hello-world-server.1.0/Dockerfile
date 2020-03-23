FROM ubuntu:16.04

WORKDIR /app

COPY helloworldserver .
COPY proxy .

EXPOSE 8080

ENV PROXY_IP=127.0.0.1
# ENV PROXY_PORT=8080
ENV SERVICE_PORT=4201

CMD ./proxy $SERVICE_PORT hello-world-server $PROXY_IP:8080 > proxy.log & sleep 5 && ./helloworldserver > server.log
