package prereqs_test

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"

	"testing"
)

func TestPrereqs(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Prereqs Suite")
}
