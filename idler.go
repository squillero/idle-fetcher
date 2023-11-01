// Idle Fetcher -- Yet another tiny Go experiment
// Copyright 2023 Giovanni Squillero
// SPDX-License-Identifier: 0BSD

package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

func timeout(out chan IpInfo, defaultValue string, duration time.Duration) {
	time.Sleep(duration)
	out <- IpInfo{
		Ip:        defaultValue,
		Source:    "Timeout",
		Timestamp: time.Now(),
	}
}

func idler() NetworkInterface {
	chanLocal := make(chan IpInfo)
	chanPublic := make(chan IpInfo)

	go timeout(chanLocal, "", 1*time.Second)
	go timeout(chanPublic, "", 1*time.Second)

	if !NoCache {
		go readCache(chanLocal, chanPublic)
	}
	go getLocalIpUDP(chanLocal)
	go getLocalIpIFACE(chanLocal)
	go getPublicIpRemote(chanPublic, "http://ipinfo.io/ip")
	go getPublicIpRemote(chanPublic, "http://ipecho.net/plain")

	info := NetworkInterface{
		LocalAddress:  <-chanLocal,
		PublicAddress: <-chanPublic,
	}

	if !NoCache {
		updateCache(info)
	}

	return info
}

func updateCache(info NetworkInterface) {
	file, _ := json.MarshalIndent(info, "", "    ")
	if err := os.WriteFile(CacheFile, file, 0600); err != nil && Verbose {
		log.Println("Error updating cache", err)
	}
}

func readCache(local chan IpInfo, public chan IpInfo) {
	var info NetworkInterface
	data, err := os.ReadFile(CacheFile)
	if err != nil {
		if Verbose {
			log.Println("Can't read cache:", err)
		}
		return
	} else {
		if Verbose {
			log.Printf("Reading '%s'\n", CacheFile)
		}
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return
	}
	diff := time.Now().Sub(info.PublicAddress.Timestamp)
	if diff < 8*time.Hour {
		info.PublicAddress.cached = true
		public <- info.PublicAddress
	}
}
