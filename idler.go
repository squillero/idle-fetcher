// Idle Fetcher -- Yet another tiny Go experiment
// Copyright 2023 Giovanni Squillero
// SPDX-License-Identifier: 0BSD

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func timeout(out chan IpInfo, defaultValue string, duration time.Duration) {
	time.Sleep(duration)
	out <- IpInfo{
		Ip:        defaultValue,
		Source:    "Timeout",
		cached:    false,
		Timestamp: time.Now(),
	}
}

func idler() NetworkInterface {
	chanLocal := make(chan IpInfo)
	chanPublic := make(chan IpInfo)

	go timeout(chanLocal, "", 1*time.Second)
	go timeout(chanPublic, "", 1*time.Second)

	go getLocalIp(chanLocal)
	go readCache(chanLocal, chanPublic)
	go getPublicIpRemote(chanPublic, "http://ipinfo.io/ip")
	go getPublicIpRemote(chanPublic, "http://ipecho.net/plain")

	info := NetworkInterface{
		LocalAddress:  <-chanLocal,
		PublicAddress: <-chanPublic,
	}

	if cache, err := os.Create(CacheFile); err == nil {
		updateCache(cache, info)
		defer cache.Close()
	}

	return info
}

func updateCache(cache *os.File, info NetworkInterface) {

	file, _ := json.MarshalIndent(info, "", "    ")
	if err := os.WriteFile(CacheFile, file, 0644); err == nil && Verbose {
		fmt.Printf("Updating cachefile '%s'\n", CacheFile)
	}

}

func readCache(local chan IpInfo, public chan IpInfo) {
	var info NetworkInterface
	if data, err := os.ReadFile(CacheFile); err == nil {
		json.Unmarshal(data, &info)
		info.LocalAddress.cached = true
		info.PublicAddress.cached = true
		public <- info.PublicAddress
	}
}
