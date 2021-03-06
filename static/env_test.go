package static

import (
	"context"

	"github.com/icecave/honeycomb/backend"
	"github.com/icecave/honeycomb/name"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("FromEnv", func() {
	Describe("fromEnv", func() {
		DescribeTable(
			"it produces the correct route",
			func(env string, expected *backend.Endpoint) {
				locator, err := fromEnv([]string{env})
				Expect(err).ShouldNot(HaveOccurred())

				endpoint, _ := locator.Locate(context.Background(), name.Parse("foo.com"))
				Expect(endpoint).To(Equal(expected))
			},
			Entry("TLS (https)", "ROUTE_FOO=foo.* https://foo.backend.com:1234", &backend.Endpoint{
				Description: "FOO",
				Address:     "foo.backend.com:1234",
				TLSMode:     backend.TLSEnabled,
			}),
			Entry("non-TLS (http)", "ROUTE_FOO=foo.* http://foo.backend.com:1234", &backend.Endpoint{
				Description: "FOO",
				Address:     "foo.backend.com:1234",
				TLSMode:     backend.TLSDisabled,
			}),
			Entry("TLS (wss)", "ROUTE_FOO=foo.* wss://foo.backend.com:1234", &backend.Endpoint{
				Description: "FOO",
				Address:     "foo.backend.com:1234",
				TLSMode:     backend.TLSEnabled,
			}),
			Entry("non-TLS (ws)", "ROUTE_FOO=foo.* ws://foo.backend.com:1234", &backend.Endpoint{
				Description: "FOO",
				Address:     "foo.backend.com:1234",
				TLSMode:     backend.TLSDisabled,
			}),
			Entry("TLS (https, implicit port)", "ROUTE_FOO=foo.* https://foo.backend.com", &backend.Endpoint{
				Description: "FOO",
				Address:     "foo.backend.com:443",
				TLSMode:     backend.TLSEnabled,
			}),
			Entry("non-TLS (http, implicit port)", "ROUTE_FOO=foo.* http://foo.backend.com", &backend.Endpoint{
				Description: "FOO",
				Address:     "foo.backend.com:80",
				TLSMode:     backend.TLSDisabled,
			}),
			Entry("TLS (wss, implicit port)", "ROUTE_FOO=foo.* wss://foo.backend.com", &backend.Endpoint{
				Description: "FOO",
				Address:     "foo.backend.com:443",
				TLSMode:     backend.TLSEnabled,
			}),
			Entry("non-TLS (ws, implicit port)", "ROUTE_FOO=foo.* ws://foo.backend.com", &backend.Endpoint{
				Description: "FOO",
				Address:     "foo.backend.com:80",
				TLSMode:     backend.TLSDisabled,
			}),
			Entry("custom description", "ROUTE_FOO=foo.* https://foo.backend.com:1234 This is the description!", &backend.Endpoint{
				Description: "This is the description!",
				Address:     "foo.backend.com:1234",
				TLSMode:     backend.TLSEnabled,
			}),
		)

		It("allows multiple routes", func() {
			env := []string{
				"ROUTE_FOO=foo.* https://foo.backend.com:1234",
				"ROUTE_BAR=bar.* https://bar.backend.com:1234",
			}

			locator, err := fromEnv(env)

			Expect(err).ShouldNot(HaveOccurred())

			endpoint, score := locator.Locate(
				context.Background(),
				name.Parse("foo.com"),
			)
			Expect(score).To(BeNumerically(">", 0))
			Expect(endpoint.Address).To(Equal("foo.backend.com:1234"))

			endpoint, score = locator.Locate(
				context.Background(),
				name.Parse("bar.com"),
			)
			Expect(score).To(BeNumerically(">", 0))
			Expect(endpoint.Address).To(Equal("bar.backend.com:1234"))
		})

		It("ignores other environment variables", func() {
			env := []string{"PATH=/usr/local/bin"}

			locator, err := fromEnv(env)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(locator).To(HaveLen(0))
		})

		It("returns an error if the match pattern is invalid", func() {
			env := []string{"ROUTE_FOO=/ https://backend"}

			_, err := fromEnv(env)

			Expect(err).Should(HaveOccurred())
		})

		It("returns an error if the URL can not be parsed", func() {
			env := []string{"ROUTE_FOO=www ://backend"}

			_, err := fromEnv(env)

			Expect(err).Should(HaveOccurred())
		})
	})
})
