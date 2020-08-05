package main

import (
	"fmt"
	"os"
)

// Config : All configuration values needed for the rendezvous api service
type Config struct {
	organizationID string
	serviceID      string
	pulsarURL      string
}

// LoadConfig : Load configuration from environment variables
func LoadConfig() Config {
	return Config{
		organizationID: os.Getenv("ORGANIZATION_ID"),
		serviceID:      os.Getenv("SERVICE_ID"),
		pulsarURL:      os.Getenv("PULSAR_URL"),
	}
}

// TopicName : Construct a topic name with the specified suffix
func (config Config) TopicName(suffix string) string {
	return fmt.Sprintf("%s/%s/%s", config.organizationID, config.serviceID, suffix)
}
