package main

import (
	"fmt"
	"log"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
)

var wg *wgctrl.Client
var cfg Config

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func main() {
	var err error
	cfg = ParseConfig("config.yml")
	checkInterval := time.Duration(cfg.Settings.Interval) * time.Second

	// TODO: For each record in the config, get the current dns A-record

	wg, err = wgctrl.New()
	if err != nil {
		log.Fatalf("Unable to create client: %e\n", err)
	}
	check()
	for range time.Tick(checkInterval) {
		check()
	}
}

func check() {
	connections, err := wg.Devices()
	if err != nil {
		log.Fatalf("Unable to get connections: %e\n", err)
	}
	for _, con := range connections {
		for _, peers := range con.Peers {
			key := peers.PublicKey.String()
			record := GetRecordWithKey(&cfg, key)
			if record == nil {
				continue
			}
			endpoint := peers.Endpoint
			if endpoint != nil {
				ip := endpoint.IP
				detectChange(record, ip.String())
			} else {
				fmt.Println("No endpoint in connection")
			}
		}
	}
}

func detectChange(record *Record, ip string) {
	if record.LastIP != ip {
		fmt.Printf("New IP detected! %s -> %s\n", record.LastIP, ip)
		record.LastIP = ip
	}
}
