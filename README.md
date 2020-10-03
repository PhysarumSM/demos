# Demos
Example services that interact with the system.

You will need the various Multi-Tier-Cloud components running in your network (allocator, registry-service, Prometheus and ping-monitor). You will need a registry-cli built on your machine that you use to build a demo and add it to the system. To interact with the demos you'll need a proxy built as well. To use registry-cli and proxy remember to either set the P2P_BOOTSTRAPS environment variable, or the --bootstrap flag. You will need a DockerHub account and will have to create DockerHub repos to host the Docker images of these demo microservices.

[Hello World](#hello-world)

[Calculator](#calculator)

[Covid19 Tracker](#covid19-tracker)

[CPU Usage Predictor](#cpu-usage-predictor)


# Hello World
Service that returns "Hello, World".

Build it.
```
$ cd helloworld/helloworldserver
$ go build
```

Add to registry-service.
```
$ <path/to/registry-cli> add service-conf.json <DockerHub repo> hello-world-server:1.1
```

To interact with it:

In one terminal start the proxy.
```
$ <path/to/proxy> 4200
```

In another terminal build and run the helloworldclient.
```
$ cd helloworld/helloworldclient
$ go build
$ ./helloworldclient 4200
```

Helloworldclient will query helloworldserver and return its response.

## Calculator
Service that performs basic arithmetic.

Build it.
```
$ cd calculator/calc-server
$ go build
```

Add to registry-service.
```
$ <path/to/registry-cli> add service-conf.json <DockerHub repo> calc-server:1.1
```

To interact with it:

In one terminal start the proxy.
```
$ <path/to/proxy> 4200
```

In another terminal build and run the calc-client.
```
$ cd calculator/calc-client
$ go build
$ ./calc-client 4200
```

Calc-client will generate a random expression every5 seconds and query calc-server for the answer.

# Covid19 Tracker
Service that tracks covid-19 statistics. Consist of a simple "database", a CLI for querying the database, and a webserver which provides a frontend and query the database in the backend.

Build the database.
```
$ cd covid19-tracker/covid19-db
$ go build
```

Add to registry-service.
```
$ <path/to/registry-cli> add service-conf.json <DockerHub repo> covid19-db:1.1
```

To interact with it:

In one terminal start the proxy.
```
$ <path/to/proxy> 4200
```

In another terminal build and run the covid19-cli.
```
$ cd covid19-tracker/covid19-cli
$ go build
$ ./covid19-cli <options> 4200
```

Options include:
```
  -city string
        City
  -country string
        Country or region
  -day string
        Day (1-31)
  -month string
        Month (1-12)
  -province string
        Province or state
  -year string
        Year (2020)
```

To run the webserver:

Note you don't need the covid19-cli from the previous step to run the webserver, but you can reuse the proxy.

Build and run the covid19-webserver.
```
$ cd covid19-tracker/covid19-webserver
$ go build
$ ./covid19-webserver 4200 8080
```

You can now visit the website hosted on you machine at port 8080 from your web browser.

Alternatively, you can use the registry-cli to build an image. This way you don't need to manually run your own proxy instance (so you won't need the proxy from the previous step). You don't need the registry-cli to push the image to DockerHub or add it to the registry-service. We can do this by using the `--no-add` flag. In this case, the DockerHub repo argument doesn't really matter, but it will still serve as the generated Docker image's name. The service name argument name doesn't really matter either. A service-conf.json is provided in the the covid19-webserver's directory.
```
$ cd covid19-tracker/covid19-webserver
$ go build
$ <path/to/registry-cli> add --no-add service-conf.json covid19-webserver covid19-webserver
```

Now run a Docker container. Make sure you set the container's P2P_BOOTSTRAPS environment variable in this command.
```
# docker container run --network="host" --detach -e "P2P_BOOTSTRAPS=<bootstrap-p2p-addr>" -e "PROXY_PORT=4200" -e "SERVICE_PORT=8080" covid19-webserver
```

You can now visit the website hosted on you machine at port 8080 from your web browser.

# CPU Usage Predictor

## Sensor
Pushes CPU readings to aggregator every 0.5 seconds. Also starts an http server that responds to all requests with "OK". This allows you to check if it is alive, and to automically allocate it using a proxy and allocator. Each sensor randomly creates an ID which is uses to send requests to aggregator. When first started, sends a request to aggregator to "prefetch" it, and then waits 10 seconds before sending data. This is to make sure an aggregator is allocated by the time sensor starts sending data, otherwise you'll get a flood of requests (one every 0.5 seconds) and a flood of aggregator instances.

To build and add to registry-service:
```
$ cd demos/cpu-usage/sensor
$ go build
$ registry-cli add service-conf.json <DockerHub repo> cpu-usage-sensor:1.0
```

To run manually, run proxy:

`$ ./proxy <proxy port> cpu-usage-sensor:1.0 <service port> <metrics port>`

Followed by the sensor:

`$ ./sensor <proxy port> <listen port>`


## Aggregator
Receives data from sensors, formats that data so it is suitable for training, and pushes to predictor. Accepts requests to endpoint `/upload/<id>/<data point>` where id is some number that indentifies the sender and data point is a number representing CPU utilization percentage. Requests to any other endpoint is responded to with "OK". Every 15 seconds, checks all received data. For each id, splits its data points into groups of 5, if there are at least data points. Then sends this formatted data to predictor. When first started, sends a request to predictor to "prefetch" it. This is to make sure a predictor is allocated by the time aggregator starts sending data, as predictor can take some time to start up.

To build and add to registry-service:
```
$ cd demos/cpu-usage/aggregator
$ go build
$ registry-cli add service-conf.json <DockerHub repo> cpu-usage-aggregator:1.0
```

## Predictor
Trains a couple ML models on data received from aggregator, plus a set of initially provided sample data. Responds to 2 endpoints. POST request a json array of arrays where the inner arrays should contain 5 numbers (eg. [[1,2,3,4,5],[6,7,8,9,10]]) to /upload. GET request to /data to return all collected data. Repeatedly performs a set of training runs on 2 linear regression models and a polynomial regression model. As new data comes in, future training runs will use the new data.

To build and add to registry-service:
```
$ cd demos/cpu-usage/predictor
$ go build
$ registry-cli add service-conf.json <DockerHub repo> cpu-usage-predictor:1.0
```

## Setup
You can trigger a bunch of sensors to be allocated by creating you own proxy and sending a bunch of requests for `cpu-usage-sensor:1.0` which should allocate them somewhere. eg. `$ ./proxy 4200` followed by `$ curl http://127.0.0.1:4200/cpu-usage-sensor:1.0`. Or boot them manually, whether by running docker containers or the manual instructions under [Sensor](#sensor). Once sensors are created, they should push data to aggregator which pushes to predictor, which should automatically allocate them.
