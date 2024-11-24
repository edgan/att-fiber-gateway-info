# Usage

```
att-fiber-gateway-info 1.0.12

Usage:
  -action string
        Action to perform (broadband-status, device-list, fiber-status, home-network-status,
                           ip-allocation, nat-check, nat-connections, nat-destinations,
                           nat-sources, nat-totals, reset-connection, reset-device,
                           reset-firewall, reset-ip, reset-wifi, restart-gateway,
                           system-information)
  -allmetrics bool (default: false)
        Return all metrics
  -cookiefile string (default: /var/tmp/att-fiber-gateway-info_cookies.gob)
        File to save session cookies
  -datadog bool (default: false)
        Send metrics to datadog
  -debug bool (default: false)
        Enable debug mode
  -filter string
        Filter to perform (icmp, ipv4, ipv6, tcp, udp)
  -fresh bool (default: false) 
        Do not use existing cookies (Warning: If always used the gateway will
        run out of sessions.)
  -interval int (default: 0)
        How often to repeat metrics
  -metrics bool (default: false)
        Return metrics based on the data instead the data
  -no bool (default: false)
        Answer no to any questions
  -password string 
        Gateway password
  -pretty bool (default: false)
        Enable pretty mode for nat-connections
  -statsdipport string
        Statsd ip port, like 127.0.0.1:8125
  -url string
        Gateway base URL
  -version bool (default: false)
        Show version
  -yes bool (default: false)
        Answer yes to any questions

```
