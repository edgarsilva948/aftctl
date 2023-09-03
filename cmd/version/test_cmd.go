/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package version

import (
	"bytes"
	"fmt"
	"io"
	"os"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = ginkgo.Describe("Validate version command", func() {
	ginkgo.Context("print the latest version of the CLI (Major, Minor, Patch and goVersion)", func() {
		ginkgo.When("the command aftctl version is executed", func() {
			ginkgo.It("should print an object with the version", func() {
				cmd := &cobra.Command{
					Use: "version",
				}

				// Redirect stdout
				old := os.Stdout // keep backup of the real stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				Run(cmd, []string{})

				// Read stdout
				outC := make(chan string)
				go func() {
					var buf bytes.Buffer
					io.Copy(&buf, r)
					outC <- buf.String()
				}()

				// Reset stdout
				w.Close()
				os.Stdout = old

				// Check stdout
				out := <-outC
				expectedRegex := `Version: {Major:"\d+", Minor:"\d+", Patch:"\d+", GoVersion:"go[\d\.]+"}`

				fmt.Println("Expected Regex:", expectedRegex)
				fmt.Println("Actual         :", out)

				gomega.Expect(out).To(gomega.MatchRegexp(expectedRegex))
			})
		})
	})
})
