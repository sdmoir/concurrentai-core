module github.com/concurrent-ai/rendezvous/src/model-enricher

go 1.14

require (
	github.com/apache/pulsar-client-go v0.1.1
	github.com/concurrent-ai/rendezvous/src/shared v0.0.0-00010101000000-000000000000
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/stretchr/testify v1.4.0
)

replace github.com/concurrent-ai/rendezvous/src/shared => ../shared