# Demos
Example services that interact with the system.

You will need the various Multi-Tier-Cloud components running in your network (allocator, registry-service, Prometheus and ping-monitor). You will need a registry-cli built on your machine that you use to build a demo and add it to the system. To interact with the demos you'll need a proxy built as well. To use registry-cli and proxy remember to either set the P2P_BOOTSTRAPS environment variable, or the --bootstrap flag. You will need a DockerHub account and will have to create DockerHub repos to host the Docker images of these demo microservices.

## Hello World
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

## Covid19 Tracker
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

## CPU Usage Predictor

Sensor. Pushes CPU readings to aggregator.
```
$ cd demos/cpu-usage/sensor
$ registry-cli add service-conf.json <DockerHub repo> cpu-usage-sensor:1.0
```
Aggregator. Formats CPU data and pushes to predictor.
```
$ cd demos/cpu-usage/aggregator
$ registry-cli add service-conf.json <DockerHub repo> cpu-usage-aggregator:1.0
```

Predictor. Trains a couple ML models on data. Responds to 2 endpoints, POST request a json array of 5 float arrays (eg. [[1,2,3,4,5],[6,7,8,9,10]]) to /upload, and GET request to /data to return all collected data.
```
$ cd demos/cpu-usage/predictor
$ registry-cli add service-conf.json <DockerHub repo> cpu-usage-predictor:1.0
```

You'll probably want to set up by triggering a bunch of sensors to be allocated. Create you own proxy and send a bunch of requests for cpu-usage-sensor:1.0 which should allocate them somewhere. Or boot them manually. From there data should get pushed to aggregator which pushes to predictor, which automatically allocates them.