# Docker
[Here](/Dockerfile) is a Dockerfile for running this in a Linux container.

## Building the image
```
docker build . -t att-fiber-gateway-info
```

# DockHub reepository
https://hub.docker.com/repository/docker/edgan/att-fiber-gateway-info/general

## Note
**The password is always required.**

## Styles
There are two styles of usage, daemon, and  command line.

Daemon uses `docker/att-fiber-gateway-info.sh` to pass in values for flags. It
also assumes you want `-allmetrics -continuous -datadog` as default flags.

Command line lets you just directly specify all the flags directly.

### Daemon environment variables for flag values
```
PASSWORD
STATSDIPPORT
URL
```

### Daemon usage
```
docker run -it -e PASSWORD='<password>' edgan/att-fiber-gateway-info:1.0.15
docker run -it -e PASSWORD='<password>' -e STATSDIPPORT='<ip>:<port>' edgan/att-fiber-gateway-info:1.0.15
docker run -it -e PASSWORD='<password>' -e URL='<url>' edgan/att-fiber-gateway-info:1.0.15
```

### Daemon examples
```
docker run -it -e PASSWORD='1234567890' edgan/att-fiber-gateway-info:1.0.15
docker run -it -e PASSWORD='1234567890' -e STATSDIPPORT='192.168.10.10:8125' edgan/att-fiber-gateway-info:1.0.15
docker run -it -e PASSWORD='1234567890' -e URL='https://192.168.10.1' edgan/att-fiber-gateway-info:1.0.15
```

### Command line usage
```
docker run -it att-fiber-gateway-info edgan/att-fiber-gateway-info:1.0.15 -action fiber-status -password "<password>"
docker run -it att-fiber-gateway-info edgan/att-fiber-gateway-info:1.0.15 -action fiber-status -metrics -password "<password>"
docker run -it att-fiber-gateway-info edgan/att-fiber-gateway-info:1.0.15 -action nat-connections -pretty  -password "<password>"
```

### Command line examples
```
docker run -it att-fiber-gateway-info edgan/att-fiber-gateway-info:1.0.15 -action fiber-status -password "1234567890"
docker run -it att-fiber-gateway-info edgan/att-fiber-gateway-info:1.0.15 -action fiber-status -metrics -password "1234567890"
docker run -it att-fiber-gateway-info edgan/att-fiber-gateway-info:1.0.15 -action nat-connections -pretty -password "1234567890"
```
