package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UpdateData struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

// GetRecord will update the last ip variable for a record based on what is already in the DNS
func GetRecord(record *Record) {
	url := fmt.Sprintf("https://cloudflare.com/client/v4/zones/%s/dns_records/%s", record.Zone, record.Record)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error making request for %s (%s): %e\n", record.Name, url, err)
	}
	req.Header.Set("X-Auth-Email", record.Email)
	req.Header.Set("X-Auth-Key", record.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error getting DNS record for %s: %e\n", record.Name, err)
		return
	}
	defer resp.Body.Close()
	// TODO: Return the DNS content
	record.LastIP = record.Name
}

// UpdateRecord will update the provided record with the provided ip in the DNS
func UpdateRecord(record *Record, ip string) {
	log.Printf("Updating Cloudflare DNS for %s to %s\n", record.Name, ip)

	data := UpdateData{
		Type: "A",
		Name: record.Name,
		Content: ip,
		TTL: record.TTL,
		Proxied: record.Proxied,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Uanble to marshal data when updating DNS record: %e\n", err)
	}
	body := bytes.NewReader(payloadBytes)

	url := fmt.Sprintf("https://cloudflare.com/client/v4/zones/%s/dns_records/%s", record.Zone, record.Record)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Fatalf("Error making request for %s (%s): %e\n", record.Name, url, err)
	}
	req.Header.Set("X-Auth-Email", record.Email)
	req.Header.Set("X-Auth-Key", record.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error updating DNS record for %s: %e\n", record.Name, err)
		return
	}
	defer resp.Body.Close()
	record.LastIP = ip
}
