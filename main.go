// Idle Fetcher -- Yet another tiny Go experiment
// Copyright 2023 Giovanni Squillero
// SPDX-License-Identifier: 0BSD

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var CacheFile = filepath.Join(os.TempDir(), "idle-fetcher.json")
var Verbose bool

type IpInfo struct {
	Ip        string
	Source    string
	Timestamp time.Time
	cached    bool
}

type NetworkInterface struct {
	LocalAddress, PublicAddress IpInfo
}

func main() {
	flag.BoolVar(&Verbose, "v", false, "Be verbose during operations")
	clearFlag := flag.Bool("c", false, "Clear cache before running")
	flag.Parse()

	if *clearFlag {
		if Verbose {
			fmt.Printf("Clearing cachefile '%s'\n", CacheFile)
		}
		os.Remove(CacheFile)
	}

	info := idler()
	fmt.Printf("%s/%s\n", info.LocalAddress.Ip, info.PublicAddress.Ip)
}
