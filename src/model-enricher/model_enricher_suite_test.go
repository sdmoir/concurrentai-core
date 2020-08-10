package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModelEnricher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ModelEnricher Suite")
}
