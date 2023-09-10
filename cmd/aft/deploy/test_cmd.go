/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package deploy contains tests for the prereqs cmd
package deploy

import (
	ginkgo "github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("testing the prereqs steps", func() {

	ginkgo.Context("stores the bucket name that will store the tfstate", func() {
		ginkgo.When("the command aftctl deploy prereqs --terraform-state-bucket-name=\"\" ", func() {
			ginkgo.It("should print an object with the version", func() {

			})
		})
	})

})
