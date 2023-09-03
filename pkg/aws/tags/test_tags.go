package tags

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Tags", func() {

	ginkgo.Context("When defining prefix", func() {
		ginkgo.It("should have a prefix as 'created-by-'", func() {
			gomega.Expect(prefix).To(gomega.Equal("created-by-"))
		})
	})

	ginkgo.Context("When using Aftctl tag", func() {
		ginkgo.It("should concatenate prefix and 'aftctl'", func() {
			gomega.Expect(Aftctl).To(gomega.Equal("created-by-aftctl"))
		})
	})

	ginkgo.Context("When using True constant", func() {
		ginkgo.It("should be 'true'", func() {
			gomega.Expect(True).To(gomega.Equal("true"))
		})
	})

})
