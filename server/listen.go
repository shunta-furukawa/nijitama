package main

import (
	"fmt"
	"net"
	"strconv"
)

// LAN mode only listens on a specific private IPv4 interface. Local development
// retains its loopback default and explicit NIJI_ADDR configuration.
func listenAddress(lanIP, port, developmentAddress string) (string, error) {
	if lanIP == "" {
		return developmentAddress, nil
	}
	ip := net.ParseIP(lanIP)
	if ip.To4() == nil || !ip.IsPrivate() {
		return "", fmt.Errorf("NIJI_LAN_IP requires a private IPv4 address")
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 {
		return "", fmt.Errorf("NIJI_PORT must be 1–65535")
	}
	return net.JoinHostPort(ip.String(), strconv.Itoa(n)), nil
}
