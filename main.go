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
	"path/filepath"
	"time"
)

const cache_file string = "idle-fetcher.cache"

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
	cache, err := os.Open(filepath.Join(os.TempDir(), cache_file))
	if err != nil {
		return
	}
	defer cache.Close()
	var local, public string
	fmt.Fscanf(cache, "%s %s", &local, &public)
	fmt.Println("ZZZZ: ", local, public)
	out <- "CACHED:" + public
}

func main() {
	chanLocal := make(chan string)
	chanPublic := make(chan string)

	go timeout(chanLocal, 1*time.Second)
	go timeout(chanPublic, 1*time.Second)

	go getLocalIp(chanLocal)
	go getPublicIpCache(chanPublic)
	go getPublicIpRemote(chanPublic, "http://ipinfo.io/ip")
	go getPublicIpRemote(chanPublic, "http://ipecho.net/plain")

	local := <-chanLocal
	public := <-chanPublic

	cache, err := os.Create(filepath.Join(os.TempDir(), cache_file))
	if err != nil {
		return
	}
	defer cache.Close()
	cache.WriteString(fmt.Sprintf("%s %s", local, public))

	fmt.Printf("%s/%s\n", local, public)
	fmt.Printf("Cache: %s\n", filepath.Join(os.TempDir(), cache_file))

}
