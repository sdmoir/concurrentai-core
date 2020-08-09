module github.com/concurrent-ai/rendezvous/src/api

go 1.14

require (
	github.com/apache/pulsar-client-go v0.1.1
	github.com/concurrent-ai/rendezvous/src/shared v0.0.0-00010101000000-000000000000
	github.com/confluentinc/confluent-kafka-go v1.4.2
	github.com/google/uuid v1.1.1
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	golang.org/x/sys v0.0.0-20200808120158-1030fc2bf1d9 // indirect
)

replace github.com/concurrent-ai/rendezvous/src/shared => ../shared
