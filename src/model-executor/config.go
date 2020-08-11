package main

import (
	"fmt"
	"os"
)

// Config : All configuration values needed for the rendezvous api service
type Config struct {
	OrganizationID string
	ServiceID      string
	ModelID        string
	ModelEndpoint  string
	PulsarURL      string
}

// LoadConfig : Load configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		OrganizationID: os.Getenv("ORGANIZATION_ID"),
		ServiceID:      os.Getenv("SERVICE_ID"),
		ModelID:        os.Getenv("MODEL_ID"),
		ModelEndpoint:  os.Getenv("MODEL_ENDPOINT"),
		PulsarURL:      os.Getenv("PULSAR_URL"),
	}
}

// TopicName : Construct a topic name with the specified suffix
func (config *Config) TopicName(suffix string) string {
	return fmt.Sprintf("%s/%s/%s", config.OrganizationID, config.ServiceID, suffix)
}
