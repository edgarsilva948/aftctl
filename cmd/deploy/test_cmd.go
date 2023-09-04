package deploy

import (
	"github.com/edgarsilva948/aftctl/cmd/deploy/aft"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Deploy Command", func() {

	ginkgo.Context("Cmd", func() {
		ginkgo.It("should have 'deploy' as the use field", func() {
			gomega.Expect(Cmd.Use).To(gomega.Equal("deploy"))
		})

		ginkgo.It("should have a slice of aliases containing 'setup'", func() {
			gomega.Expect(Cmd.Aliases).To(gomega.ConsistOf("setup"))
		})

		ginkgo.It("should have a short description", func() {
			gomega.Expect(Cmd.Short).To(gomega.Equal("Deploy AFT from from stdin"))
		})

		ginkgo.It("should have a long description", func() {
			gomega.Expect(Cmd.Long).To(gomega.Equal("Deploy AFT from from stdin"))
		})

		ginkgo.It("should register aft.Cmd as a subcommand", func() {
			gomega.Expect(Cmd.HasSubCommands()).To(gomega.BeTrue())
			gomega.Expect(Cmd.Commands()).To(gomega.ContainElement(aft.Cmd))
		})
	})

})
