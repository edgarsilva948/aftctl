/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// This file contains test information about the tool.

package info

import (
	"bytes"
	"fmt"
	"runtime"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Information Utilities", func() {

	ginkgo.Describe("Fetching Go Runtime Version", func() {
		ginkgo.It("should accurately return the Go runtime version", func() {
			gomega.Expect(GetGoVersion()).To(gomega.Equal(runtime.Version()))
		})
	})

	ginkgo.Describe("Assembling Current Tool Version", func() {
		var v Version

		ginkgo.BeforeEach(func() {
			v = BuildCurrentVersion()
		})

		ginkgo.Context("when assembling version details", func() {
			ginkgo.It("should properly format the Major version", func() {
				gomega.Expect(v.Major).To(gomega.MatchRegexp(`[0-9]+`))
			})

			ginkgo.It("should properly format the Minor version", func() {
				gomega.Expect(v.Minor).To(gomega.MatchRegexp(`[0-9]+`))
			})

			ginkgo.It("should properly format the Patch version", func() {
				gomega.Expect(v.Patch).To(gomega.MatchRegexp(`[0-9]+`))
			})

			ginkgo.It("should match the system's Go version", func() {
				gomega.Expect(v.GoVersion).To(gomega.Equal(runtime.Version()))
			})
		})
	})

	ginkgo.Describe("Displaying Version Info", func() {
		ginkgo.It("matches expected output", func() {
			var buf bytes.Buffer
			PrintVersion(&buf)

			v := BuildCurrentVersion()
			expectedOutput := fmt.Sprintf("Version: {Major:\"%s\", Minor:\"%s\", Patch:\"%s\", GoVersion:\"%s\"}\n",
				v.Major, v.Minor, v.Patch, v.GoVersion)
			gomega.Expect(buf.String()).To(gomega.Equal(expectedOutput))
		})
	})
})
