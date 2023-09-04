/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package color

import (
	"errors"
	"os"
	"runtime"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

type mockOsInfo struct {
	goos           string
	stdoutStatErr  error
	stdoutStatMode os.FileMode
}

func (m mockOsInfo) GetOs() string {
	return m.goos
}

func (m mockOsInfo) GetStdoutStat() (os.FileInfo, error) {
	return nil, m.stdoutStatErr
}

var _ = ginkgo.Describe("Enabling the Color flag", func() {

	var cmd *cobra.Command
	var args []string
	var toComplete string

	ginkgo.BeforeEach(func() {
		cmd = &cobra.Command{}
		args = []string{}
		toComplete = ""
	})

	ginkgo.Context("UseColor function", func() {

		ginkgo.When("getting the real OS", func() {
			ginkgo.It("should return a string with the OS", func() {
				realOsInfo := RealOsInfo{}
				expectedOS := runtime.GOOS

				gomega.Expect(realOsInfo.GetOs()).To(gomega.Equal(expectedOS))
			})
		})

		ginkgo.When("GetStdoutStat function", func() {
			ginkgo.It("should return a valid FileInfo and no error", func() {
				realOsInfo := RealOsInfo{}

				fileInfo, err := realOsInfo.GetStdoutStat()

				var _ os.FileInfo = fileInfo

				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("completion function", func() {

			ginkgo.It("should return valid completion options and directive", func() {
				options, directive := completion(cmd, args, toComplete)
				gomega.Expect(options).To(gomega.ConsistOf("auto", "never", "always"))
				gomega.Expect(directive).To(gomega.Equal(cobra.ShellCompDirectiveDefault))
			})
		})

		ginkgo.When("when color option is 'never'", func() {
			ginkgo.It("should return false", func() {
				color = "never"
				osInfo := mockOsInfo{goos: "linux"}
				gomega.Expect(UseColor(osInfo)).To(gomega.BeFalse())
			})
		})

		ginkgo.When("when color option is 'always'", func() {
			ginkgo.It("should return true", func() {
				color = "always"
				osInfo := mockOsInfo{goos: "linux"}
				gomega.Expect(UseColor(osInfo)).To(gomega.BeTrue())
			})
		})

		ginkgo.When("AddFlag function", func() {
			var cmd *cobra.Command
			ginkgo.BeforeEach(func() {
				cmd = &cobra.Command{}
				AddFlag(cmd)
			})

			ginkgo.It("should add a 'color' flag", func() {
				gomega.Expect(cmd.PersistentFlags().Lookup("color")).To(gomega.Not(gomega.BeNil()))
			})
		})

		ginkgo.When("when runtime is windows", func() {
			ginkgo.It("should return false", func() {
				color = "auto"
				osInfo := mockOsInfo{goos: "windows"}
				gomega.Expect(UseColor(osInfo)).To(gomega.BeFalse())
			})
		})

		ginkgo.When("when os.Stdout.Stat returns an error", func() {
			ginkgo.It("should return true", func() {
				color = "auto"
				osInfo := mockOsInfo{goos: "linux", stdoutStatErr: errors.New("some error")}
				gomega.Expect(UseColor(osInfo)).To(gomega.BeFalse())
			})
		})

		ginkgo.When("when stdout.Mode has specific flags", func() {
			ginkgo.It("should return based on flag conditions", func() {
				color = "auto"
				osInfo := mockOsInfo{goos: "linux", stdoutStatMode: os.ModeDevice}
				expected := (osInfo.stdoutStatMode&os.ModeDevice == 0) && (osInfo.stdoutStatMode&os.ModeNamedPipe == 0)
				gomega.Expect(UseColor(osInfo)).To(gomega.Equal(expected))
			})
		})
	})
})
