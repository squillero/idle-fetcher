// Idle Fetcher -- Yet another tiny Go experiment
// Copyright 2023 Giovanni Squillero
// SPDX-License-Identifier: 0BSD

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var CacheFile = filepath.Join(os.TempDir(), "idle-fetcher.json")

// FLAGS
var Verbose bool = false
var NoCache bool = false

type IpInfo struct {
	Ip        string
	Source    string
	Timestamp time.Time
	reliable  bool
	cached    bool
}

type NetworkInterface struct {
	LocalAddress, PublicAddress IpInfo
}

func main() {
	flag.BoolVar(&Verbose, "v", false, "Verbose operations")
	flag.BoolVar(&NoCache, "n", false, "Don't use cache")
	clearCache := flag.Bool("r", false, "Remove cache before running")
	flag.Parse()

	log.SetPrefix("[IdleFetcher] ")

	if *clearCache {
		if err := os.Remove(CacheFile); err == nil && Verbose {
			log.Printf("Removed '%s'\n", CacheFile)
		}
	}

	info := idler()

	if Verbose {
		var cached string
		if info.LocalAddress.cached {
			cached = " (cached)"
		} else {
			cached = ""
		}
		log.Printf("Got local IP info from %s%s\n", info.LocalAddress.Source, cached)
		if info.PublicAddress.cached {
			cached = " (cached)"
		} else {
			cached = ""
		}
		log.Printf("Got public IP info from %s%s\n", info.PublicAddress.Source, cached)
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "hal9001"
	}

	if info.LocalAddress.Ip == "" {
		fmt.Printf("%s\n", hostname)
	} else if info.PublicAddress.Ip == "" || info.LocalAddress.Ip == info.PublicAddress.Ip {
		fmt.Printf("%s [%s]\n", hostname, info.LocalAddress.Ip)
	} else {
		fmt.Printf("%s [%s/%s]\n", hostname, info.LocalAddress.Ip, info.PublicAddress.Ip)
	}
}
