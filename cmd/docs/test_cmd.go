/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package docs

import (
	"os"
	"path/filepath"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = ginkgo.Describe("Cobra docs package to generate documentation", func() {
	ginkgo.Context("the function CreateDir creates the output dir if it doesn exists", func() {
		ginkgo.When("the command aftctl docs is executed", func() {
			ginkgo.It("should create a directory if it does not exist", func() {
				dir, err := os.MkdirTemp("", "example")
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				nonExistentDir := filepath.Join(dir, "nonexistent")

				err = CreateDir(nonExistentDir)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				_, err = os.Stat(nonExistentDir)
				gomega.Expect(os.IsNotExist(err)).To(gomega.BeFalse())
			})

			ginkgo.It("should not return an error if the directory already exists", func() {
				dir, err := os.MkdirTemp("", "example")
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				err = CreateDir(dir)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})
	})
})

var _ = ginkgo.Describe("Cobra docs package to generate documentation", func() {
	ginkgo.Context("when format is markdown", func() {
		ginkgo.When("the command aftctl docs is executed", func() {
			ginkgo.It("should create markdown documentation", func() {
				dir, err := os.MkdirTemp("", "docstest")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(dir)

				args.dir = dir
				args.format = "markdown"

				cmd := &cobra.Command{Use: "docs"}
				err = Run(cmd, []string{})
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Context("when format is man", func() {
		ginkgo.When("the command aftctl docs is executed", func() {
			ginkgo.It("should create man documentation", func() {
				dir, err := os.MkdirTemp("", "docstest")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(dir)

				args.dir = dir
				args.format = "man"

				cmd := &cobra.Command{Use: "docs"}
				err = Run(cmd, []string{})
				gomega.Expect(err).To(gomega.BeNil())

			})
		})
	})
})
