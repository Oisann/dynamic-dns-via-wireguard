package main

import (
	"flag"
	"log"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
)

var wg *wgctrl.Client
var cfg Config
var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.yml", "Specify config file")
	flag.Parse()
	cfg = ParseConfig(configFile)

	for i, _ := range cfg.Records {
		GetRecord(&cfg.Records[i])
	}

	DisplayConfig(&cfg)
	var err error
	wg, err = wgctrl.New()
	if err != nil {
		log.Fatalf("Unable to create client: %e\n", err)
	}
}

func main() {
	check()
	checkInterval := time.Duration(cfg.Settings.Interval) * time.Second
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
