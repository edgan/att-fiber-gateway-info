# Docker

## Note
**The password is a required.**

## Styles
There are two styles of usage, command line, and daemon.

Daemon uses `docker/att-fiber-gateway-info.sh` to pass in values for flags. It
also assumes you want `-allmetrics -continuous -datadog` as default flags.

Command line lets you just directly specify all the flags directly.

## Daemon environment variables for flag values
```
PASSWORD
STATSDIPPORT
URL
```

## Building the image
```
docker build . -t att-fiber-gateway-info
```

## Command line usage
```
docker run -it att-fiber-gateway-info att-fiber-gateway-info -action fiber-status -password "<password>"
docker run -it att-fiber-gateway-info att-fiber-gateway-info -action fiber-status -metrics -password "<password>"
docker run -it att-fiber-gateway-info att-fiber-gateway-info -action nat-connections -pretty  -password "<password>"
```

## Command line examples
```
docker run -it att-fiber-gateway-info att-fiber-gateway-info -action fiber-status -password "1234567890"
docker run -it att-fiber-gateway-info att-fiber-gateway-info -action fiber-status -metrics -password "1234567890"
docker run -it att-fiber-gateway-info att-fiber-gateway-info -action nat-connections -pretty -password "1234567890"
```

### Daemon usage
```
docker run -it -e PASSWORD='<password>' att-fiber-gateway-info
docker run -it -e PASSWORD='<password>' -e STATSDIPPORT='<ip>:<port>' att-fiber-gateway-info
docker run -it -e PASSWORD='<password>' -e URL='<url>' att-fiber-gateway-info
```

### Daemon examples
```
docker run -it -e PASSWORD='1234567890' att-fiber-gateway-info
docker run -it -e PASSWORD='1234567890' -e STATSDIPPORT='192.168.10.10:8125' att-fiber-gateway-info
docker run -it -e PASSWORD='1234567890' -e URL='https://192.168.10.1' att-fiber-gateway-info
```
