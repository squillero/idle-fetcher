// Idle Fetcher -- Yet another tiny Go experiment
// Copyright 2023 Giovanni Squillero
// SPDX-License-Identifier: 0BSD

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func httpGet(url string) string {
	result, err := http.Get(url)
	if err != nil {
		return ""
	}
	cooked, err := io.ReadAll(result.Body)
	if err != nil {
		return ""
	}
	fmt.Println(string(cooked))
	return string(cooked)
}

func main() {
	fmt.Println(os.TempDir())

	go httpGet("http://ipinfo.io/ip")
	go httpGet("http://ipecho.net/plain")
	//fmt.Println("\"%s\"", httpGet("http://ipinfo.io/ip"))
	//fmt.Println("\"%s\"", httpGet("http://ipecho.net/plain"))
	time.Sleep(1 * time.Second)
}
