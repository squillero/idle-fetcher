// Idle Fetcher -- Yet another tiny Go experiment
// Copyright 2023 Giovanni Squillero
// SPDX-License-Identifier: 0BSD

package main

import (
	"net"
	"os"
	"time"
)

func getLocalIpUDP(out chan IpInfo) {
	if conn, err := net.Dial("udp", "8.8.8.8:80"); err == nil {
		defer conn.Close()
		localAddress := conn.LocalAddr().(*net.UDPAddr)
		out <- IpInfo{
			Ip:        localAddress.IP.String(),
			Source:    "UDP/conn.LocalAddr",
			Timestamp: time.Now(),
			reliable:  true,
		}
	}
}

func getLocalIpIFACE(out chan IpInfo) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				out <- IpInfo{
					Ip:        ipnet.IP.String(),
					Source:    "net.IPNet",
					Timestamp: time.Now(),
					reliable:  true,
				}
			}
		}
	}
}

func getLocalHostname(out chan IpInfo) {
	time.Sleep(1800 * time.Millisecond)
	if hostname, err := os.Hostname(); err == nil {
		out <- IpInfo{
			Ip:        hostname,
			Source:    "os.Hostname",
			Timestamp: time.Now(),
		}
	}
}
