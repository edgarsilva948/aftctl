package color

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = ginkgo.Describe("Color", func() {

	ginkgo.Describe("UseColor function", func() {
		ginkgo.Context("when color option is 'never'", func() {
			ginkgo.It("should return false", func() {
				color = "never"
				gomega.Expect(UseColor()).To(gomega.BeFalse())
			})
		})

		ginkgo.Context("when color option is 'always'", func() {
			ginkgo.It("should return true", func() {
				color = "always"
				gomega.Expect(UseColor()).To(gomega.BeTrue())
			})
		})
	})

	ginkgo.Describe("AddFlag function", func() {
		var cmd *cobra.Command
		ginkgo.BeforeEach(func() {
			cmd = &cobra.Command{}
			AddFlag(cmd)
		})

		ginkgo.It("should add a 'color' flag", func() {
			gomega.Expect(cmd.PersistentFlags().Lookup("color")).To(gomega.Not(gomega.BeNil()))
		})
	})
})
