package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Settings Settings
	Records []Record `yaml:"records"`
}
type Settings struct {
	Interval int `yaml:"interval"`
}
type Record struct {
	Name string `yaml:"name"`
	Key string `yaml:"key"`
	Token string `yaml:"token"`
	Email string `yaml:"email"`
	TTL int `yaml:"ttl"`
	Zone string `yaml:"zone"`
	Record string `yaml:"record"`
	LastIP string
}

// ParseConfig takes a file path as input and outputs a config struct
func ParseConfig(filename string) Config {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error handling config file named %s: %e\n", filename, err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Unable to decode config: %e\n", err)
	}
	return cfg
}

// GetRecordWithKey returns the Record with the specified key
func GetRecordWithKey(config *Config, key string) *Record {
	for _, record := range config.Records {
		if record.Key == key {
			return &record
		}
	}
	return nil
}
