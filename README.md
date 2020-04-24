[![<ORG_NAME>](https://circleci.com/gh/kennedyoliveira/prometheus-msi-afterburner-exporter.svg?style=shield)](https://app.circleci.com/pipelines/github/kennedyoliveira/prometheus-msi-afterburner-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/kennedyoliveira/prometheus-msi-afterburner-exporter)](https://goreportcard.com/report/github.com/kennedyoliveira/prometheus-msi-afterburner-exporter)
[![Docker Pulls](https://img.shields.io/docker/pulls/kennedyoliveira/afterburner-exporter.svg?cacheSeconds=3600)](https://hub.docker.com/r/kennedyoliveira/afterburner-exporter)

# Prometheus Exporter for MSI Afterburner
Export metrics from MSI Afterburner, any metric, to prometheus.

## Pre requisites
You need [MSI Afterburner Remote Server](http://download.msi.com/uti_exe/vga/MSIAfterburnerRemoteServer.zip) running 
on the computer with MSI Afterburner, this application allow the metrics to be queried via an API.

You can click on the above link or check [in MSI Afterburner Page](https://www.msi.com/page/afterburner) for the download.
Just download it, unzip and run it.

The `afterburner-exporter` doesn't need to be on the same computer, it can be in any computer, although it is very lightweight and
barely consume any resource, so if you run it alongside your gaming computer it won't affect performance.

## Configuration Parameters
You can pass these parameters as arguments:

Example:
```shell script
afterburner-exporter host=192.168.1.32 port=1082 listen-address=0.0.0.0:9090 metrics-endpoint=/custom/metrics
```

| Parameter         | Default               | Description
| ---------         | -------               | ------------
| host              | 127.0.0.1             | The host of the computer running MSI Afterburner Server
| port              | 82                    | The port of the MSI Afterburner Server
| username          | MSIAfterburner        | Username to authenticate in MSI Afterburner Server, should be MSIAfterburner unless a new version changes it.
| password          | 17cc95b4017d496f82    | Password to authenticate in MSI Afterburner Server, it is fixed unless you change it in the config files.
| listen-address    | 0.0.0.0:8080          | Address and port where this app will listen to request.
| metrics-endpoint  | /metrics              | Endpoint which the metrics will be available to be scrapped by Prometheus.

## Docker images
You can also use it via docker:

```shell script
docker container run --name afterburner-exporter -p 8080:8080 kennedyoliveira/afterburner-exporter
``` 

There are images built for the following architectures:

    - linux/amd64
    - linux/i386
    - linux/arm64
    - linux/arm/v7
    
Given the arm archs, you can use it on a raspberry pi or similar arm hardware.

## Build
To build just clone the project and run:
```shell script
make build
```

To cross compile for the different platforms use:
```shell script
make compile
```

For other additional options check the `Makefile`

## TODO

 - [ ] Allow blacklist from config
 - [ ] Allow gpu regex from config
 - [ ] Configuration from environment variables 