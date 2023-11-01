// Idle Fetcher -- Yet another tiny Go experiment
// Copyright 2023 Giovanni Squillero
// SPDX-License-Identifier: 0BSD

package main

import (
	"net"
	"time"
)

func getLocalIp(out chan IpInfo) {
	if conn, err := net.Dial("udp", "8.8.8.8:80"); err == nil {
		defer conn.Close()
		localAddress := conn.LocalAddr().(*net.UDPAddr)
		out <- IpInfo{
			Ip:        localAddress.IP.String(),
			Source:    "UDP",
			cached:    true,
			Timestamp: time.Now(),
		}
	}
}
