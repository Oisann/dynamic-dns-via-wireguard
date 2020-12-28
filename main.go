package main

import (
	"fmt"
	"log"

	"golang.zx2c4.com/wireguard/wgctrl"
)

func main() {
	wg, err := wgctrl.New()
	if err != nil {
		log.Fatalf("Unable to create client: %e", err)
	}
	connections, err := wg.Devices()
	if err != nil {
		log.Fatalf("Unable to get connections: %e", err)
	}
	for _, con := range connections {
		for _, peers := range con.Peers {
			fmt.Println(peers.Endpoint)
		}
	}
}
