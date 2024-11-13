# att-fiber-gateway-info
## Description
A [golang](https://en.wikipedia.org/wiki/Go_(programming_language)) command line
tool to pull values from the pages of an
[AT&T Fiber](https://www.att.com/internet/fiber/)
[BGW320-505 gateway](https://help.sonic.com/hc/en-us/articles/1500000066642-BGW320).

## Compiling
```
go build
```

## Usage
```
./att-fiber-gateway-info --help
Usage of ./att-fiber-gateway-info:
  -action string
    	Action to perform (broadband-status, fiber-status, nat-check, nat-connections, nat-destinations, nat-sources, nat-totals, sys-info)
  -cookiefile string
    	File to save session cookies (default "/var/tmp/.att-fiber-gateway-info_cookies.gob")
  -debug
    	Enable debug mode
  -filter string
    	Filter to perform (icmp, ipv4, ipv6, tcp, udp)
  -fresh
    	Do not use existing cookies (Warning: If you use all the time you will run out of sessions. There is a max.)
  -password string
    	Gateway password
  -url string
    	Gateway base URL (default "https://192.168.1.254")

```

## Actions
### broadband-status
This return all the values from the [Broadband|Status](https://192.168.1.254/cgi-bin/broadbandstatistics.ha) page.

Example:
```
Broadband Connection Source: FIBER

Broadband Connection: Up
Broadband Network Type: Lightspeed
Broadband IPv4 Address: 110.112.23.57 
Gateway IPv4 Address: 110.112.23.1
MAC Address: 0c:7c:28:10:2c:e9
Primary DNS: 68.94.156.1
Secondary DNS: 68.94.157.1
Primary DNS Name: 
Secondary DNS Name: 
MTU: 1500
Status: Available
Service Type: native IPv6
Global Unicast IPv6 Address: 2001:506:6041:a283::
Link Local Address: fe80::e7c:28ff:feb5:bddc
Default IPv6 Gateway Address: fe80::45:28:69:0
Primary DNS: 
Secondary DNS: 
MTU: 1500
Receive Packets: 14739661
Transmit Packets: 6200161
Receive Bytes: 20566453127
Transmit Bytes: 794551431
Receive Unicast: 14739651
Transmit Unicast: 6199895
Receive Multicast: 10
Transmit Multicast: 184
Receive Drops: 0
Transmit Drops: 14
Receive Errors: 0
Transmit Errors: 0
Collisions: 0
Transmit Packets: 548
Transmit Errors: 0
Transmit Discards: 0
PON Link Status: OPERATION (O5)
UNI Status: up
```

### fiber-status
This return all the values from the [Broadband|Fiber status](https://192.168.1.254/cgi-bin/fiberstat.ha).

Example:
```
Optical WAN Operational Status: Up
Fiber Module: Unavailable
Last Change: 1731528578
Link State: Up
Name: SFP
Connector: 1
Transceiver: 200000000000000000
Encoding: 3
BR Nominal: 100
Br Min: 0
Br Max: 0
Rate ID: 100
Wave Length: 1270 nm
Tx Disable State: 0
RS1 State: 0
Rate Select State: 0
Tx Fault State: 0
Rx LOS State: 0
Data Ready Bar State: 0

Length SMF-km: 40
Length SMF: 0
Length 50uM: 0
Length 62dot5uM: 0
Length OM3: 0

Vendor Name: NOKIA
Vendor OUI: 000000
Vendor PN: 3FE46901-STCC
Vendor Rev: 1B
Vendor SN: DBT31020200199
Vendor Date Code: 240401

OPT Cooled Trans: uncooled transceiver
OPT Powerlvl: 1
OPT Linear Rcvr: conventional receiver
OPT Rate Select: 0
OPT Tx Disable: 1
OPT Tx Fault: 1
OPT Inverted-LOS: 0
OPT LOS: 1

DMC Type Legacy: 0
DMC Type Implemented: 1
DMC Type Internal Cal: 1
DMC Type External Cal: 0
DMC Type Rx Avg Pwr: Average power method

EOC Alarm Implemented: 1
EOC Soft Tx Disable: 1
EOC Soft Tx fault: 1
EOC Soft Rx LOS: 1
EOC Soft Rate Select: 0

SFF 8079 App Select: 0
SFF 8431 Soft Rate Select: 0
SFF Ver Compliance: rev 11.0

Temperature: 44.76
Vcc: 3.38
Tx Bias: 111.28
Tx Power: 39.692
Rx Power: 0.54
```

### nat-check
This returns the total number of connections as a `float`. I am using this mode
as a metric to be monitored via [Datadog](https://www.datadoghq.com/). In my
case I will get an alert via [PagerDuty](https://www.pagerduty.com/) if the
value goes over **8000**.

Example:
```
45.0
```

### nat-connections
This returns all the values from the table on the
[Diagnostics|NAT Table](https://192.168.1.254/cgi-bin/nattable.ha) page.

Example:
```
IP Family Protocol Protocol Number Lifetime TCP State   Source Address                          Source Port Destination Address                     Destination Port NAT Source Address                      NAT Source Port NAT Destination Address                 NAT Destination Port Bidirectional 
ipv4      udp      17              18                   110.112.23.57                          48996       8.8.8.8                                 53               110.112.23.57                          48996           8.8.8.8                                 53                   *             
ipv4      udp      17              441                  110.112.23.57                          48727       8.8.8.8                                 53               110.112.23.57                          48727           8.8.8.8                                 53                   *
```

### nat-destinations
This counts the destination addresses occurrences from the table, and sorts the
results in descending order.

Example:
```
Destinations IP addresses:
6, 192.168.1.254
5, 110.112.23.57
4, 8.8.8.8
3, 142.250.141.188
1, 3.233.148.21
1, 161.38.184.18
1, 173.194.8.198
1, 142.250.68.42
1, 35.186.224.31
1, 3.233.147.249
1, 172.217.14.106
1, 142.251.2.188
```

### nat-sources
This counts the source addresses occurrences from the table, and sorts the
results in descending order.

Example:
```
Source IP addresses:
151, 110.112.23.57
5, 192.168.254.254
2, 162.125.40.1
2, 142.250.189.14
2, 89.134.5.161
```

### nat-totals
This reports the total number of connections, 
[tcp](https://en.wikipedia.org/wiki/Transmission_Control_Protocol) connections,
and [udp](https://en.wikipedia.org/wiki/User_Datagram_Protocol) connections.
The [icmp](https://en.wikipedia.org/wiki/Internet_Control_Message_Protocol)
connections are not listed, but are part of the total.

Example:
```
Total number of connections: 143
Total number of tcp connections: 95
Total number of udp connections: 43
```

### sys-info
This returns the values from the [Device|System Information](https://192.168.1.254/cgi-bin/sysinfo.ha) page.

Example:
```
Manufacturer: NOKIA
Model Number: BGW320-505
Serial Number: N94VA0JP111713
Software Version: 6.28.7
MAC Address: 0c:7c:28:10:2c:e9
First Use Date: 2020/11/16 17:19:53
Time Since Last Reboot: 00:09:51:01
Current Date/Time: 2024-11-13T13:13:31
Hardware Version: 02001E0046004F
```

## Known issues
1. By default assumes `/var/tmp` exists. This can be worked around with the
`-cookiefile` command line argument. It is untested on
[Windows](https://en.wikipedia.org/wiki/Microsoft_Windows).

## Story
  I just had [AT&T Fiber](https://www.att.com/internet/fiber/) installed. As
part of the installation process I received a
[AT&T](https://www.att.com/)
([Nokia](https://www.nokia.com/))
[BGW320-505 gateway](https://help.sonic.com/hc/en-us/articles/1500000066642-BGW320).
I immediately found and set up the
[passthrough mode](https://www.devonstephens.com/how-to-enable-ip-passthrough-on-att-bgw320-505/).
I already had my own [iptables](https://en.wikipedia.org/wiki/Iptables) based
router running [Fedora Linux](https://fedoraproject.org/). I wanted to
continue to use it. I wasn't interested in disabling my existing
[DHCP](https://en.wikipedia.org/wiki/Dynamic_Host_Configuration_Protocol)
servers and [DNS](https://en.wikipedia.org/wiki/Domain_Name_System) servers. I
also wasn't interested in changing the ip addresses of my whole network.

In my research one downside I found of
[passthrough mode](https://www.devonstephens.com/how-to-enable-ip-passthrough-on-att-bgw320-505/)
is that the
[BGW320-505 gateway](https://help.sonic.com/hc/en-us/articles/1500000066642-BGW320)
continues to do some connection tracking. It has a max of **8192** connections.
This can be tracked on the
[Diagnostics|NAT Table](https://192.168.1.254/cgi-bin/nattable.ha) page. It is
said that the most common way this becomes a problem is via
[bittorrent](https://en.wikipedia.org/wiki/BitTorrent) clients opening enough
connections to go over the limit.

The problem with the
[Diagnostics|NAT Table](https://192.168.1.254/cgi-bin/nattable.ha) page is it
is behind a [login](https://192.168.1.254/cgi-bin/login.ha) page. I wanted to
write a script to scrape the page. I then dug into the
[HTML](https://en.wikipedia.org/wiki/HTML) form and
[javascript](https://en.wikipedia.org/wiki/JavaScript). What I found is the
authors don't trust the network given you can access it via
[HTTP](https://en.wikipedia.org/wiki/HTTP) or
HTTPS](https://en.wikipedia.org/wiki/HTTPS) with
an invalid certificate. So they made the login process complicate.

There seems to a bug in the page design where the first time you go it will
claim your browser doesn't accept cookies, but it works on reload. My guess
is that this happens because it is trying to read the cookie that it hasn't
given you yet.

The login process goes something like this.

1. Load a page, get redirected to the
[login](https://192.168.1.254/cgi-bin/login.ha) page, and get told your browser
doesn't accept cookies.
2. Reload the page.
3. Receive the [nonce](https://en.wikipedia.org/wiki/Cryptographic_nonce)
hidden value in the page, and hashing functions written in
[javascript](https://en.wikipedia.org/wiki/JavaScript) that use
[md5](https://en.wikipedia.org/wiki/MD5).
[md5](https://en.wikipedia.org/wiki/MD5)(password+
[nonce](https://en.wikipedia.org/wiki/Cryptographic_nonce)).
4. Posting to the [login](https://192.168.1.254/cgi-bin/login.ha) page these
variables.
  - [nonce](https://en.wikipedia.org/wiki/Cryptographic_nonce), an all lower
case 64 character hex-decimal value. Looks like a
[sha256sum](https://en.wikipedia.org/wiki/SHA-2), like
`87428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7`.
  - password, not the actual password, but instead the same number of
characters as the password replaced with `*` characters
  - hashpassword, the results of the hashpwd function
  - Continue, the submit button that calls the hashing function and triggers
the post.

Once you are logged in your session is only tracked by your session cookie. It
is possible to generate too many sessions. I accomplished this during the
writing of this code, because at first it was generating a new session per run.

### Programming languages
I originally prototyped this written in
[javascript](https://en.wikipedia.org/wiki/JavaScript). Since I was trying to
reuse the [javascript](https://en.wikipedia.org/wiki/JavaScript) from the login
page. I started with [PhantomJS](https://github.com/ariya/phantomjs). It
worked, and was using [Firefox](https://www.mozilla.org/en-US/firefox/). This
meant [Firefox](https://www.mozilla.org/en-US/firefox/) needed to be installed,
and it needed to be able to run in the background. I switched to
[Puppeteer](https://github.com/puppeteer/puppeteer) once I realized
[PhantomJS](https://github.com/ariya/phantomjs) was not being maintained.
[Puppeteer](https://github.com/puppeteer/puppeteer) worked better, but suffered
the same problems as [PhantomJS](https://github.com/ariya/phantomjs).

I then realized I could just rewrite the
[javascript](https://en.wikipedia.org/wiki/JavaScript) code in
[golang](https://en.wikipedia.org/wiki/Go_(programming_language)). This works
well. It greatly simplifies the dependecies, CPU resources, memory resources,
testing, and more.

## Flaws in the login page
There are two flaws in the login process.

1. The security of the whole process is tied to the security of the session
cookie. Given that it goes across
[HTTP](https://en.wikipedia.org/wiki/HTTP)(no encryption) or
[HTTPS](https://en.wikipedia.org/wiki/HTTPS) with an invalid
certificate(easy to man in the middle), this is not good.
2. The hashing method is just a fig leaf, because of flaw one. In addition it
uses [md5](https://en.wikipedia.org/wiki/MD5), which is known to be insecure.
Modern solutions would be things like [sha256](https://en.wikipedia.org/wiki/SHA-2)
or [bcrypt](https://en.wikipedia.org/wiki/Bcrypt).

As far as I can tell there is no way to replace the gateway's
[SSL](https://en.wikipedia.org/wiki/Transport_Layer_Security) certificate.

I first ran across the first flaw over 15 years ago. I had gone to a
[Defcon](https://en.wikipedia.org/wiki/DEF_CON) talk about it,
and then discovered the same flaw in the real world.
