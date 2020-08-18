module github.com/concurrentai/concurrentai-core/src/rendezvous-collector

go 1.14

require (
	github.com/apache/pulsar-client-go v0.1.1
	github.com/concurrentai/concurrentai-core/src/shared v0.0.0-00010101000000-000000000000
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
)

replace github.com/concurrentai/concurrentai-core/src/shared => ../shared
