# Actions
## broadband-status
This return all the values from the
[Broadband|Status](https://192.168.1.254/cgi-bin/broadbandstatistics.ha) page.

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

## device-list
This returns all the values from the
[Device|Device List](https://192.168.1.254/cgi-bin/devices.ha) page.

Example:
```
MAC Address: f4:f5:d8:c8:10:39
IPv4 Address: 192.168.1.102
Name: Google-Home
Last Activity: Thu Nov 14 16:48:47 2024
Status: on
Allocation: dhcp
Connection Type:
  Ethernet LAN-1
Connection Speed: 2500Mbps    fullduplex
Mesh Client: No
IPv6 Address: fe80::f6f5:d8ff:fec8:1039
Type: slaac
Valid Lifetime: forever
Preferred Lifetime: forever

MAC Address: fe:7d:07:99:03:10
Name: unknownfe7d07990310
Last Activity: Mon Nov 11 17:37:06 2024
Status: off
Allocation: dhcp
Connection Type:
  Wi-Fi: 2.4 GHz
  Type: Home
  Name: ATTAJXBcOw
Connection Speed
Mesh Client: No
```

## fiber-status
This returns all the values from the
[Broadband|Fiber Status](https://192.168.1.254/cgi-bin/fiberstat.ha) page.

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

## home-network-status
This returns the values from the
[Home Network|Status](https://192.168.1.254/cgi-bin/lanstatistics.ha) page.

Example:
```
Device IPv4 Address:    192.168.1.254
DHCPv4 Netmask:         255.255.255.0
DHCP Server:            On
DHCPv4 Start Address:   192.168.1.100
DHCPv4 End Address:     192.168.1.199
DHCP Leases Available:  81
DHCP Leases Allocated:  19
DHCP Primary Pool:      Private
Secondary Subnet:       Disabled
Public Subnet:
Cascaded Router Status: Disabled
IP Passthrough Status:  Off (private IP address)

Interface:     Status   Active Devices Inactive Devices
Ethernet:      Enabled  0              4
5G Ethernet:   Enabled  18             9
Wi-Fi 2.4 GHz: Disabled 0              1
Wi-Fi 5 GHz:   Disabled 0              1
Mesh Clients:  Disabled 0              0

Status: Unavailable

Transmit Packets:  1133082
Transmit Errors:   0
Transmit Discards: 0
Receive Packets:   395522
Receive Errors:    0
Receive Discards:  2

                                                        2.4 GHz  5 GHz
Wi-Fi Radio Status:                                     Disabled Disabled
Wi-Fi is not enabled. Click here to configure Wi-Fi on.

                    Port 1     Port 2 Port 3 Port 4
State:              up         down   down   down
Transmit Speed:     2500000000 0      0      0
Transmit Packets:   1135624    0      0      0
Transmit Bytes:     993216374  0      0      0
Transmit Unicast:   767318     0      0      0
Transmit Multicast: 159980     0      0      0
Transmit Dropped:   0          0      0      0
Transmit Errors:    0          0      0      0
Receive Packets:    396330     0      0      0
Receive Bytes:      111567069  0      0      0
Receive Unicast:    342931     0      0      0
Receive Multicast:  37871      0      0      0
Receive Dropped:    0          0      0      0
Receive Errors:     0          0      0      0

```

## ip-allocation
This returns the values from the
[Home Network|IP Allocation](https://192.168.1.254/cgi-bin/ipalloc.ha) page.

Example:
```
IPv4 Address / Name                        MAC Address        Status  Allocation       Action
minipc-1                                   00:1e:06:48:2f:a9  off     DHCP Allocation
lake                                       00:1f:c6:fc:35:91  off     DHCP Allocation
192.168.1.137 / ASUSTek COMPUTER INC.      00:23:54:1c:26:55  on      DHCP Allocation
unknown0024e4f453a6                        00:24:e4:f4:53:a6  off     DHCP Allocation
192.168.1.138 / river                      00:c0:ca:13:7f:51  on      DHCP Allocation
192.168.1.154 / android-dhcp-13            0a:e2:62:1c:fd:b5  on      DHCP Allocation
192.168.1.143 / BRW2C98113D0FE5            2c:98:11:3d:0F:e5  on      DHCP Allocation
192.168.1.147 / fuchsia-1cf2-9a48-aba0     3c:8d:20:e3:3d:7b  on      DHCP Allocation
192.168.1.161 / My-ecobee                  44:61:32:f8:fa:2d  on      DHCP Allocation
router                                     48:21:0b:6a:5f:51  off     DHCP Allocation
192.168.1.107 / Pixel-8-Pro                5a:aa:40:d9:35:17  on      DHCP Allocation
192.168.1.101 / Vickies-MBP                88:66:5a:4f:93:46  on      DHCP Allocation
192.168.1.159 / ocean                      a0:36:bc:1c:3f:92  on      DHCP Allocation
192.168.1.106 / VY-Pixel-6                 d2:33:fa:5c:d3:61  on      DHCP Allocation
Pixel-8-Pro                                d6:7d:1a:1b:e8:3f  off     DHCP Allocation
192.168.1.117 / Vickies-5th-Gen-iPad-Pro   d6:b7:73:c2:c3:f3  on      DHCP Allocation
192.168.1.135 / higgs                      d6:f5:4a:1c:a3:1a  off     DHCP Allocation
192.168.1.136 / VY-Pixel-6                 de:4a:ff:cd:f5:13  off     DHCP Allocation
192.168.1.103 / Google-Home                f4:f5:d8:4f:de:1a  on      DHCP Allocation
192.168.1.104 / Google-Home                f4:f5:d8:b2:c1:31  on      DHCP Allocation
192.168.1.102 / Google-Home                f4:f5:d8:c8:10:39  on      DHCP Allocation
unknownfe7d07990310                        fe:7d:07:99:03:10  off     DHCP Allocation
```

## nat-check
This returns the total number of connections as a `float`. I am using this mode
as a metric to be monitored via [Datadog](https://www.datadoghq.com/). In my
case I will get an alert via [PagerDuty](https://www.pagerduty.com/) if the
value goes over **8000**.

Example:
```
45.0
```

## nat-connections
This returns all the values from the table on the
[Diagnostics|NAT Table](https://192.168.1.254/cgi-bin/nattable.ha) page.

It has a default comma delimited mode, and a pretty mode with aligned columns.

Example default:
```
IP Family, Protocol, Protocol Number, Lifetime, TCP State, Source Address, Source Port, Destination Address, Destination Port, NAT Source Address, NAT Source Port, NAT Destination Address, NAT Destination Port, Bidirectional
ipv4, udp, 17, 28, , 192.168.1.103, 5353, 192.168.1.254, 5353, 192.168.1.103, 5353, 192.168.1.254, 5353,
ipv4, udp, 17, 463, , 192.168.1.104, 33049, 8.8.8.8, 53, 110.112.23.57, 33049, 8.8.8.8, 53, *
```

Example pretty:
```
IP Family  Protocol  Protocol Number  Lifetime  TCP State    Source Address                           Source Port  Destination Address                      Destination Port  NAT Source Address                       NAT Source Port  NAT Destination Address                  NAT Destination Port  Bidirectional
ipv4       udp       17               28                     192.168.1.103                            5353         192.168.1.254                            5353              192.168.1.103                            5353             192.168.1.254                            5353
ipv4       udp       17               485                    192.168.1.106                            57970        64.233.177.95                            443               110.112.23.57                            57970            64.233.177.95                            443                   *
```

## nat-destinations
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

## nat-sources
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

## nat-totals
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

## reset-connection
This asks you if you are sure you want to reset the connection. If you answer
yes it resets the router via
[Dianostics|Resets](https://192.168.1.254/cgi-bin/resets.ha) page. There is
also the -yes flag to skip the question.

Be careful with this action.

## reset-device
This asks you if you are sure you want to reset the device. If you answer
yes it resets the router via
[Dianostics|Resets](https://192.168.1.254/cgi-bin/resets.ha) page. There is
also the -yes flag to skip the question.

Be careful with this action.

## reset-firewall
This asks you if you are sure you want to reset the firewall. If you answer
yes it resets the router via
[Dianostics|Resets](https://192.168.1.254/cgi-bin/resets.ha) page. There is
also the -yes flag to skip the question.

Be careful with this action.

## reset-ip
This asks you if you are sure you want to reset the ip. If you answer
yes it resets the router via
[Dianostics|Resets](https://192.168.1.254/cgi-bin/resets.ha) page. There is
also the -yes flag to skip the question.

Be careful with this action.

## reset-wifi
This asks you if you are sure you want to reset the wifi. If you answer
yes it resets the router via
[Dianostics|Resets](https://192.168.1.254/cgi-bin/resets.ha) page. There is
also the -yes flag to skip the question.

Be careful with this action.

## restart-gateway
This asks you if you are sure you want to restart the gateway. If you answer
yes it restarts the router via
[Dianostics|Resets](https://192.168.1.254/cgi-bin/resets.ha) page. There is
also the -yes flag to skip the question.

Be careful with this action.

## system-information
This returns the values from the
[Device|System Information](https://192.168.1.254/cgi-bin/sysinfo.ha) page.

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
