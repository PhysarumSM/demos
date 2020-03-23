FROM ubuntu:16.04

WORKDIR /app

COPY helloworldserver .
COPY proxy .

EXPOSE 8080

ENV PROXY_SERVICE_IP=127.0.0.1

CMD ./proxy hello-world-server $PROXY_SERVICE_IP:8080 > proxy.log & sleep 5 && ./helloworldserver > server.log
