# Metrics

There is the `-metrics` flag that returns the metrics for the action specified. It looks for tables with `Statistics` at the end of the name.

There is also the `-allmetrics` flag that returns mertics for all actions known to have metrics.

## Naming and formatting
It pulls the `model` from the `System Information` page returned by the `system-information` action. It converts dashes and spaces to dots. It adds `.0` the end to make it a `float` for reporting to [Datadog](https://www.datadoghq.com/) as metrics.


```
model.action.metric=value
            |
            V
bgw320-505.broadband-status.Receive Packets=46538166
            |
            V
bgw320505.broadband.status.receive.packets=46538166.0
```

## Usage
```
att-fiber-gateway-info -action broadband-status -metrics
att-fiber-gateway-info -action fiber-status -metrics
att-fiber-gateway-info -action home-network-status -metrics
att-fiber-gateway-info -allmetrics
```

## Examples

### broadband-status
```
bgw320505.broadband.status.receive.packets=46538166.0
bgw320505.broadband.status.transmit.packets=11006967.0
bgw320505.broadband.status.receive.bytes=1133156935.0
bgw320505.broadband.status.transmit.bytes=1160606147.0
bgw320505.broadband.status.receive.unicast=46537844.0
bgw320505.broadband.status.transmit.unicast=11006967.0
bgw320505.broadband.status.receive.multicast=323.0
bgw320505.broadband.status.transmit.multicast=323.0
bgw320505.broadband.status.receive.drops=0.0
bgw320505.broadband.status.transmit.drops=7.0
bgw320505.broadband.status.receive.errors=0.0
bgw320505.broadband.status.transmit.errors=0.0
bgw320505.broadband.status.collisions=0.0
bgw320505.broadband.status.transmit.packets=11006967.0
bgw320505.broadband.status.transmit.errors=0.0
bgw320505.broadband.status.transmit.discards=7.0
```

### fiber-status
```
bgw320505.fiber.status.temperature=42
bgw320505.fiber.status.vcc=3
bgw320505.fiber.status.tx.bias=11
bgw320505.fiber.status.tx.power=59
bgw320505.fiber.status.rx.power=-223
```

### home-network-status
```
bgw320505.home.network.status.transmit.packets=45229281.0
bgw320505.home.network.status.transmit.errors=0.0
bgw320505.home.network.status.transmit.discards=0.0
bgw320505.home.network.status.receive.packets=9654441.0
bgw320505.home.network.status.receive.errors=0.0
bgw320505.home.network.status.receive.discards=0.0
bgw320505.home.network.status.port1.state=up
bgw320505.home.network.status.port2.state=down
bgw320505.home.network.status.port3.state=down
bgw320505.home.network.status.port4.state=down
bgw320505.home.network.status.port1.transmit.speed=2500000000.0
bgw320505.home.network.status.port2.transmit.speed=0.0
bgw320505.home.network.status.port3.transmit.speed=0.0
bgw320505.home.network.status.port4.transmit.speed=0.0
bgw320505.home.network.status.port1.transmit.packets=43206738.0
bgw320505.home.network.status.port2.transmit.packets=0.0
bgw320505.home.network.status.port3.transmit.packets=0.0
bgw320505.home.network.status.port4.transmit.packets=0.0
bgw320505.home.network.status.port1.transmit.bytes=708750858.0
bgw320505.home.network.status.port2.transmit.bytes=0.0
bgw320505.home.network.status.port3.transmit.bytes=0.0
bgw320505.home.network.status.port4.transmit.bytes=0.0
bgw320505.home.network.status.port1.transmit.unicast=42913167.0
bgw320505.home.network.status.port2.transmit.unicast=0.0
bgw320505.home.network.status.port3.transmit.unicast=0.0
bgw320505.home.network.status.port4.transmit.unicast=0.0
bgw320505.home.network.status.port1.transmit.multicast=128388.0
bgw320505.home.network.status.port2.transmit.multicast=0.0
bgw320505.home.network.status.port3.transmit.multicast=0.0
bgw320505.home.network.status.port4.transmit.multicast=0.0
bgw320505.home.network.status.port1.transmit.dropped=0.0
bgw320505.home.network.status.port2.transmit.dropped=0.0
bgw320505.home.network.status.port3.transmit.dropped=0.0
bgw320505.home.network.status.port4.transmit.dropped=0.0
bgw320505.home.network.status.port1.transmit.errors=0.0
bgw320505.home.network.status.port2.transmit.errors=0.0
bgw320505.home.network.status.port3.transmit.errors=0.0
bgw320505.home.network.status.port4.transmit.errors=0.0
bgw320505.home.network.status.port1.receive.packets=9150050.0
bgw320505.home.network.status.port2.receive.packets=0.0
bgw320505.home.network.status.port3.receive.packets=0.0
bgw320505.home.network.status.port4.receive.packets=0.0
bgw320505.home.network.status.port1.receive.bytes=3146498964.0
bgw320505.home.network.status.port2.receive.bytes=0.0
bgw320505.home.network.status.port3.receive.bytes=0.0
bgw320505.home.network.status.port4.receive.bytes=0.0
bgw320505.home.network.status.port1.receive.unicast=8869237.0
bgw320505.home.network.status.port2.receive.unicast=0.0
bgw320505.home.network.status.port3.receive.unicast=0.0
bgw320505.home.network.status.port4.receive.unicast=0.0
bgw320505.home.network.status.port1.receive.multicast=152219.0
bgw320505.home.network.status.port2.receive.multicast=0.0
bgw320505.home.network.status.port3.receive.multicast=0.0
bgw320505.home.network.status.port4.receive.multicast=0.0
bgw320505.home.network.status.port1.receive.dropped=159.0
bgw320505.home.network.status.port2.receive.dropped=0.0
bgw320505.home.network.status.port3.receive.dropped=0.0
bgw320505.home.network.status.port4.receive.dropped=0.0
bgw320505.home.network.status.port1.receive.errors=0.0
bgw320505.home.network.status.port2.receive.errors=0.0
bgw320505.home.network.status.port3.receive.errors=0.0
bgw320505.home.network.status.port4.receive.errors=0.0
```
