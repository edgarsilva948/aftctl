/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package local_test

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"

	"testing"
)

func TestVersion(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Local Suite")
}
