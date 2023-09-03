package completion

import (
	"bytes"
	"io"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = ginkgo.Describe("Completion", func() {

	var (
		cmd *cobra.Command
		out io.Writer
	)

	ginkgo.BeforeEach(func() {
		cmd = &cobra.Command{Use: "completion"}
		out = &bytes.Buffer{}
	})

	ginkgo.Context("RunCompletion function", func() {

		ginkgo.When("no arguments are provided", func() {
			ginkgo.It("should generate Bash completion by default", func() {
				RunCompletion(cmd, []string{}, out)
				gomega.Expect(out.(*bytes.Buffer).String()).To(gomega.ContainSubstring("# bash completion for completion"))
			})
		})

		ginkgo.When("bash argument is provided", func() {
			ginkgo.It("should generate Bash completion", func() {
				RunCompletion(cmd, []string{"bash"}, out)
				gomega.Expect(out.(*bytes.Buffer).String()).To(gomega.ContainSubstring("# bash completion for completion"))
			})
		})

		ginkgo.When("zsh argument is provided", func() {
			ginkgo.It("should generate Zsh completion", func() {
				RunCompletion(cmd, []string{"zsh"}, out)
				gomega.Expect(out.(*bytes.Buffer).String()).To(gomega.ContainSubstring("# zsh completion for completion"))
			})
		})

		ginkgo.When("powershell argument is provided", func() {
			ginkgo.It("should generate PowerShell completion", func() {
				RunCompletion(cmd, []string{"powershell"}, out)
				gomega.Expect(out.(*bytes.Buffer).String()).To(gomega.ContainSubstring("# powershell completion for completion"))
			})
		})
	})
})
