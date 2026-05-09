package design

import . "goa.design/goa/v3/dsl"

var BasicAuthPayload = Type("BasicAuthPayload", func() {
	Description("HTTP basic auth credentials.")
	Username("user", String)
	Password("pass", String)
})

var BasicAuth = BasicAuthSecurity("basic")

var _ = API("service", func() {
	Title("service")
	Description("News Deframer Goa design.")
	Version("1.0.0")
	Server("service", func() {
		Host("localhost", func() {
			URI("http://0.0.0.0:8080")
		})
	})
})
