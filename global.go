package main

// Actions is an immutable slice of strings
var actions = []string{
	"broadband-status", "device-list", "fiber-status", "home-network-status",
	"ip-allocation", "nat-check", "nat-connections", "nat-destinations",
	"nat-sources", "nat-totals", "reset-connection", "reset-device",
	"reset-firewall", "reset-ip", "reset-wifi", "restart-gateway",
	"system-information",
}

var filters = []string{"icmp", "ipv4", "ipv6", "tcp", "udp"}

var metricActions = []string{
	"broadband-status", "fiber-status", "home-network-status", "nat-totals",
}
