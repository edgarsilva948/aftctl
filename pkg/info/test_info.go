package info

import (
	"runtime"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Info", func() {

	ginkgo.Describe("GetGoVersion function", func() {
		ginkgo.It("should return the Go runtime version", func() {
			gomega.Expect(GetGoVersion()).To(gomega.Equal(runtime.Version()))
		})
	})

	ginkgo.Describe("BuildCurrentVersion function", func() {
		var v Version

		ginkgo.BeforeEach(func() {
			v = BuildCurrentVersion()
		})

		ginkgo.Context("when fetching version", func() {
			ginkgo.It("should have Major matching the pattern", func() {
				gomega.Expect(v.Major).To(gomega.MatchRegexp(`[0-9]+`))
			})

			ginkgo.It("should have Minor matching the pattern", func() {
				gomega.Expect(v.Minor).To(gomega.MatchRegexp(`[0-9]+`))
			})

			ginkgo.It("should have Patch matching the pattern", func() {
				gomega.Expect(v.Patch).To(gomega.MatchRegexp(`[0-9]+`))
			})
			ginkgo.It("should have correct GoVersion", func() {
				gomega.Expect(v.GoVersion).To(gomega.Equal(runtime.Version()))
			})
		})
	})

	// Note: Testing PrintVersion might require capturing stdout
})
