package color

import (
	"os"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

type mockOsInfo struct {
	goos string
}

func (m mockOsInfo) GetOs() string {
	return m.goos
}

func (m mockOsInfo) GetStdoutStat() (os.FileInfo, error) {
	return nil, nil
}

var _ = ginkgo.Describe("Color", func() {

	ginkgo.Describe("UseColor function", func() {
		ginkgo.Context("when color option is 'never'", func() {
			ginkgo.It("should return false", func() {
				color = "never"
				osInfo := mockOsInfo{goos: "linux"}
				gomega.Expect(UseColor(osInfo)).To(gomega.BeFalse())
			})
		})

		ginkgo.Context("when color option is 'always'", func() {
			ginkgo.It("should return true", func() {
				color = "always"
				osInfo := mockOsInfo{goos: "linux"}
				gomega.Expect(UseColor(osInfo)).To(gomega.BeTrue())
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

	ginkgo.Context("when runtime is windows", func() {
		ginkgo.It("should return false", func() {
			color = "auto"
			osInfo := mockOsInfo{goos: "windows"}
			gomega.Expect(UseColor(osInfo)).To(gomega.BeFalse())
		})
	})

})
