# Metrics

The general idea is to be able to pull metrics out of various gateway pages.

## Screenshots
![Sreenshot1](/screenshots/datadog-metrics1.png)
![Sreenshot2](/screenshots/datadog-metrics2.png)
![Sreenshot3](/screenshots/datadog-metrics3.png)

## Datadog
[Datadog](https://www.datadoghq.com/) is just one way to use this. It is free
for up to five hosts.

## Dashboards
In the [datadog/dashboards](datadog/dashboards/) folder you will find dashboard
json files for the different models of the gateway. These should be used with
`-allmetrics`.

These have the most interesting things listed first, have the time set to
global to allow you to change the amount of time, and the correct metrics
names per model.

## Usage
```
att-fiber-gateway-info -action broadband-status -metrics
att-fiber-gateway-info -action fiber-status -metrics
att-fiber-gateway-info -action home-network-status -metrics
att-fiber-gateway-info -action nat-totals -metrics
att-fiber-gateway-info -action nat-totals -metrics -interval 10
att-fiber-gateway-info -action broadband-status -metrics -datadog -interval 10
att-fiber-gateway-info -action fiber-status -metrics -datadog -interval 10
att-fiber-gateway-info -action home-network-status -metrics -datadog -interval 10
att-fiber-gateway-info -action nat-totals -metrics -datadog
att-fiber-gateway-info -allmetrics
att-fiber-gateway-info -allmetrics -datadog -interval 20
```

## Systemd
[Here](systemd/att-fiber-gateway-info.service) is a tested systemd service file
for running this as a daemon on a Linux system as root.

You should make sure there is a configuration file,
`/root/.att-fiber-gateway-info.yml`. Also that it is properly configured before
starting the service.

## Flags
There is the `-metrics` flag that returns the metrics for the action specified.
It looks for tables with certain summaries. The exception is the action
`fiber-status`.

There is also the `-allmetrics` flag that returns mertics for all actions known
to have metrics.

Finally there is the `-interval` flag to allow you to automatically return
metrics every X number of seconds. There is a minimum of 10 seconds when using
`-metrics` for individual actions, and a minimum of 20 seconds when using
`-allmetrics` to return all metrics.

## Statsd
### What is statsd?
StatsD is an industry-standard technology stack for monitoring applications and
instrumenting any piece of software to deliver custom metrics.

### Flags
There is finally the `-datadog` flag that instead of printing the metrics sends
them to statsd as configured by either `-statsdipport` or `statdIPPort` in the
configuration file. It defaults to `127.0.0.1:8125`. It only sends the `float`
metrics, because that is what `datadog` accepts for metrics. A way to deal with
`string` metrics is planned.

### Installation
Statd can be installed via statsd(Node.JS), datadog-agent(Python), or
datadog-dogstatsd(Golang).

## Naming and formatting
It pulls the `model` from the `System Information` page returned by the
`system-information` action. It converts dashes and spaces to dots. All strings
are lower cases.  It adds `.0` the end to make it a `float` for reporting to
[Datadog](https://www.datadoghq.com/) as a metric.


```
model.action.summary.metric=value
            |
            V
bgw320-505.broadband-status.IPv4.Receive Packets=46538166
            |
            V
bgw320505.broadband.status.ipv4.receive.packets=46538166.0
```

## Examples

### broadband-status
```
bgw320505.broadband.status.ethernet.line.state=up
bgw320505.broadband.status.ethernet.current.speed=10000.0
bgw320505.broadband.status.ethernet.current.duplex=full
bgw320505.broadband.status.ipv4.receive.packets=100033343.0
bgw320505.broadband.status.ipv4.transmit.packets=18178781.0
bgw320505.broadband.status.ipv4.receive.bytes=670337417.0
bgw320505.broadband.status.ipv4.transmit.bytes=3968010139.0
bgw320505.broadband.status.ipv4.receive.unicast=100032972.0
bgw320505.broadband.status.ipv4.transmit.unicast=18178781.0
bgw320505.broadband.status.ipv4.receive.multicast=371.0
bgw320505.broadband.status.ipv4.transmit.multicast=371.0
bgw320505.broadband.status.ipv4.receive.drops=0.0
bgw320505.broadband.status.ipv4.transmit.drops=7.0
bgw320505.broadband.status.ipv4.receive.errors=0.0
bgw320505.broadband.status.ipv4.transmit.errors=0.0
bgw320505.broadband.status.ipv4.collisions=0.0
bgw320505.broadband.status.ipv6.transmit.packets=18178781.0
bgw320505.broadband.status.ipv6.transmit.errors=0.0
bgw320505.broadband.status.ipv6.transmit.discards=7.0
```

### fiber-status
```
bgw320505.fiber.status.temperature=41.0
bgw320505.fiber.status.vcc=3.0
bgw320505.fiber.status.tx.bias=11.0
bgw320505.fiber.status.tx.power=56.0
bgw320505.fiber.status.rx.power=-222.0
```

### home-network-status
```
bgw320505.home.network.status.ipv4.transmit.packets=98095373.0
bgw320505.home.network.status.ipv4.transmit.errors=0.0
bgw320505.home.network.status.ipv4.transmit.discards=0.0
bgw320505.home.network.status.ipv4.receive.packets=15976962.0
bgw320505.home.network.status.ipv4.receive.errors=0.0
bgw320505.home.network.status.ipv4.receive.discards=0.0
bgw320505.home.network.status.lan.ethernet.port1.state=up
bgw320505.home.network.status.lan.ethernet.port2.state=down
bgw320505.home.network.status.lan.ethernet.port3.state=down
bgw320505.home.network.status.lan.ethernet.port4.state=down
bgw320505.home.network.status.lan.ethernet.port1.transmit.speed=2500000000.0
bgw320505.home.network.status.lan.ethernet.port2.transmit.speed=0.0
bgw320505.home.network.status.lan.ethernet.port3.transmit.speed=0.0
bgw320505.home.network.status.lan.ethernet.port4.transmit.speed=0.0
bgw320505.home.network.status.lan.ethernet.port1.transmit.packets=96111994.0
bgw320505.home.network.status.lan.ethernet.port2.transmit.packets=0.0
bgw320505.home.network.status.lan.ethernet.port3.transmit.packets=0.0
bgw320505.home.network.status.lan.ethernet.port4.transmit.packets=0.0
bgw320505.home.network.status.lan.ethernet.port1.transmit.bytes=4065353315.0
bgw320505.home.network.status.lan.ethernet.port2.transmit.bytes=0.0
bgw320505.home.network.status.lan.ethernet.port3.transmit.bytes=0.0
bgw320505.home.network.status.lan.ethernet.port4.transmit.bytes=0.0
bgw320505.home.network.status.lan.ethernet.port1.transmit.unicast=95582905.0
bgw320505.home.network.status.lan.ethernet.port2.transmit.unicast=0.0
bgw320505.home.network.status.lan.ethernet.port3.transmit.unicast=0.0
bgw320505.home.network.status.lan.ethernet.port4.transmit.unicast=0.0
bgw320505.home.network.status.lan.ethernet.port1.transmit.multicast=201313.0
bgw320505.home.network.status.lan.ethernet.port2.transmit.multicast=0.0
bgw320505.home.network.status.lan.ethernet.port3.transmit.multicast=0.0
bgw320505.home.network.status.lan.ethernet.port4.transmit.multicast=0.0
bgw320505.home.network.status.lan.ethernet.port1.transmit.dropped=0.0
bgw320505.home.network.status.lan.ethernet.port2.transmit.dropped=0.0
bgw320505.home.network.status.lan.ethernet.port3.transmit.dropped=0.0
bgw320505.home.network.status.lan.ethernet.port4.transmit.dropped=0.0
bgw320505.home.network.status.lan.ethernet.port1.transmit.errors=0.0
bgw320505.home.network.status.lan.ethernet.port2.transmit.errors=0.0
bgw320505.home.network.status.lan.ethernet.port3.transmit.errors=0.0
bgw320505.home.network.status.lan.ethernet.port4.transmit.errors=0.0
bgw320505.home.network.status.lan.ethernet.port1.receive.packets=15472570.0
bgw320505.home.network.status.lan.ethernet.port2.receive.packets=0.0
bgw320505.home.network.status.lan.ethernet.port3.receive.packets=0.0
bgw320505.home.network.status.lan.ethernet.port4.receive.packets=0.0
bgw320505.home.network.status.lan.ethernet.port1.receive.bytes=520713796.0
bgw320505.home.network.status.lan.ethernet.port2.receive.bytes=0.0
bgw320505.home.network.status.lan.ethernet.port3.receive.bytes=0.0
bgw320505.home.network.status.lan.ethernet.port4.receive.bytes=0.0
bgw320505.home.network.status.lan.ethernet.port1.receive.unicast=15189952.0
bgw320505.home.network.status.lan.ethernet.port2.receive.unicast=0.0
bgw320505.home.network.status.lan.ethernet.port3.receive.unicast=0.0
bgw320505.home.network.status.lan.ethernet.port4.receive.unicast=0.0
bgw320505.home.network.status.lan.ethernet.port1.receive.multicast=154025.0
bgw320505.home.network.status.lan.ethernet.port2.receive.multicast=0.0
bgw320505.home.network.status.lan.ethernet.port3.receive.multicast=0.0
bgw320505.home.network.status.lan.ethernet.port4.receive.multicast=0.0
bgw320505.home.network.status.lan.ethernet.port1.receive.dropped=159.0
bgw320505.home.network.status.lan.ethernet.port2.receive.dropped=0.0
bgw320505.home.network.status.lan.ethernet.port3.receive.dropped=0.0
bgw320505.home.network.status.lan.ethernet.port4.receive.dropped=0.0
bgw320505.home.network.status.lan.ethernet.port1.receive.errors=0.0
bgw320505.home.network.status.lan.ethernet.port2.receive.errors=0.0
bgw320505.home.network.status.lan.ethernet.port3.receive.errors=0.0
bgw320505.home.network.status.lan.ethernet.port4.receive.errors=0.0
```

### nat-totals
```
bgw320505.nat.totals.connections=370.0
bgw320505.nat.totals.icmp.connections=0.0
bgw320505.nat.totals.tcp.connections=178.0
bgw320505.nat.totals.udp.connections=192.0
```
