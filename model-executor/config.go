package main

import (
	"fmt"
	"os"
)

// Config : All configuration values needed for the rendezvous api service
type Config struct {
	organizationID string
	serviceID      string
	modelID        string
	pulsarURL      string
	modelEndpoint  string
}

// LoadConfig : Load configuration from environment variables
func LoadConfig() Config {
	return Config{
		organizationID: os.Getenv("ORGANIZATION_ID"),
		serviceID:      os.Getenv("SERVICE_ID"),
		modelID:        os.Getenv("MODEL_ID"),
		pulsarURL:      os.Getenv("PULSAR_URL"),
		modelEndpoint:  os.Getenv("MODEL_ENDPOINT"),
	}
}

// TopicName : Construct a topic name with the specified suffix
func (config Config) TopicName(suffix string) string {
	return fmt.Sprintf("%s/%s/%s", config.organizationID, config.serviceID, suffix)
}
