package main

import (
	"fmt"
	"log"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
)

var wg *wgctrl.Client
var lastIP string

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func main() {
	var err error

	// TODO: Check DNS record to get the current "lastIP"
	lastIP = "N/A"

	wg, err = wgctrl.New()
	if err != nil {
		log.Fatalf("Unable to create client: %e", err)
	}
	check()
	doEvery(2*time.Second, check)
}

func check() {
	connections, err := wg.Devices()
	if err != nil {
		log.Fatalf("Unable to get connections: %e", err)
	}
	for _, con := range connections {
		for _, peers := range con.Peers {
			key := peers.PublicKey.String()
			if key != "TEMTZpKBxQ26vzk/xs1lXyiwve8nQg3vhj6MlBP06l0=" {
				continue
			}
			endpoint := peers.Endpoint
			if endpoint != nil {
				ip := endpoint.IP
				detectChange(ip.String())
			} else {
				fmt.Println("No endpoint in connection")
			}
		}
	}
}

func detectChange(ip string) {
	if lastIP != ip {
		fmt.Printf("New IP detected! %s -> %s\n", lastIP, ip)
	}
	lastIP = ip
}
