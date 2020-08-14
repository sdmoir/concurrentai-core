package main

import (
	"fmt"
	"os"
)

// Config : All configuration values needed for the rendezvous api service
type Config struct {
	OrganizationID string
	ServiceID      string
	ActiveModelID  string
	PulsarURL      string
}

// LoadConfig : Load configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		OrganizationID: os.Getenv("ORGANIZATION_ID"),
		ServiceID:      os.Getenv("SERVICE_ID"),
		ActiveModelID:  os.Getenv("ACTIVE_MODEL_ID"),
		PulsarURL:      os.Getenv("PULSAR_URL"),
	}
}

// TopicName : Construct a topic name with the specified suffix
func (config *Config) TopicName(suffix string) string {
	return fmt.Sprintf("%s/%s/%s", config.OrganizationID, config.ServiceID, suffix)
}

// SubscriptionName : Construct a subscription name with the specified suffix
func (config *Config) SubscriptionName(suffix string) string {
	return fmt.Sprintf("%s/%s/%s-subscription", config.OrganizationID, config.ServiceID, suffix)
}
