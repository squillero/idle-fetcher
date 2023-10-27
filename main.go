// Idle Fetcher -- Yet another tiny Go experiment
// Copyright 2023 Giovanni Squillero
// SPDX-License-Identifier: 0BSD

package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func getLocalIp(out chan string) {
	if conn, err := net.Dial("udp", "8.8.8.8:80"); err == nil {
		defer conn.Close()
		localAddress := conn.LocalAddr().(*net.UDPAddr)
		out <- localAddress.IP.String()
	}
}

func getPublicIpRemote(out chan string, url string) {
	result, err := http.Get(url)
	if err != nil {
		return
	}
	cooked, err := io.ReadAll(result.Body)
	if err != nil {
		return
	}
	out <- url + ":" + string(cooked)
}

func getPublicIpCache(out chan string) {
	out <- "Dummy"
}

func main() {
	info := make(chan string)

	fmt.Println(getLocalIp())
	fmt.Println(os.TempDir())

	go getIpCache(info)
	go getPublicIpRemote(info, "http://ipinfo.io/ip")
	go getPublicIpRemote(info, "http://ipecho.net/plain")
	//fmt.Println("\"%s\"", getPublicIpRemote("http://ipinfo.io/ip"))
	//fmt.Println("\"%s\"", getPublicIpRemote("http://ipecho.net/plain"))

	fmt.Println(<-info)
}
