// Idle Fetcher -- Yet another tiny Go experiment
// Copyright 2023 Giovanni Squillero
// SPDX-License-Identifier: 0BSD

package main

import (
	"io"
	"net/http"
	"time"
)

func getPublicIpRemote(out chan IpInfo, url string) {
	result, err := http.Get(url)
	if err != nil {
		return
	}
	cooked, err := io.ReadAll(result.Body)
	if err != nil {
		return
	}
	out <- IpInfo{
		Ip:        string(cooked),
		Source:    url,
		Timestamp: time.Now(),
	}
}
