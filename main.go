package main

import (
	"log"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
)

var wg *wgctrl.Client
var cfg Config

func main() {
	var err error
	cfg = ParseConfig("config.yml")
	checkInterval := time.Duration(cfg.Settings.Interval) * time.Second

	for i, _ := range cfg.Records {
		GetRecord(&cfg.Records[i])
	}

	DisplayConfig(&cfg)

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
				log.Println("No endpoint in connection")
			}
		}
	}
}

func detectChange(record *Record, ip string) {
	if record.LastIP != ip {
		log.Printf("New IP detected for %s: %s -> %s\n", record.Name, record.LastIP, ip)
		go UpdateRecord(record, ip)
	}
}
