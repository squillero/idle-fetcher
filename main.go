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
	"time"
)

func timeout(out chan string, duration time.Duration) {
	time.Sleep(duration)
	out <- "?"
}

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
	local := make(chan string)
	public := make(chan string)

	go timeout(local, 1*time.Second)
	go timeout(public, 1*time.Second)

	fmt.Println(os.TempDir())

	go getLocalIp(local)
	go getPublicIpRemote(public, "http://ipinfo.io/ip")
	go getPublicIpRemote(public, "http://ipecho.net/plain")
	//fmt.Println("\"%s\"", getPublicIpRemote("http://ipinfo.io/ip"))
	//fmt.Println("\"%s\"", getPublicIpRemote("http://ipecho.net/plain"))

	fmt.Printf("%s/%s\n", <-local, <-public)
}
