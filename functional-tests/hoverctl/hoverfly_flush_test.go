package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hoverctl flush cache", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.SetMode("simulate")
		hoverfly.ImportSimulation(functional_tests.JsonPayload)
		hoverfly.Proxy(sling.New().Get("http://destination-server.com"))

		WriteConfiguration("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort())
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	It("should flush cache", func() {
		output := functional_tests.Run(hoverctlBinary, "flush", "--force")

		Expect(output).To(ContainSubstring("Successfully flushed cache"))

		cacheView := hoverfly.GetCache()

		Expect(cacheView.Cache).To(HaveLen(0))
	})

	It("should error nicely when trying to flush but cache is disabled", func() {
		hoverfly.Stop()
		hoverfly.Start("-disable-cache")
		output := functional_tests.Run(hoverctlBinary, "flush", "--force")

		Expect(output).To(ContainSubstring("Cache was not set on Hoverfly"))
	})

})
