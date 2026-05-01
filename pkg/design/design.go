package design

import . "goa.design/goa/v3/dsl"

var BasicAuthPayload = Type("BasicAuthPayload", func() {
	Username("user", String)
	Password("pass", String)
})

var BasicAuth = BasicAuthSecurity("basic")

var _ = API("service_new", func() {
	Title("service_new")
	Description("service_new Goa design")
	Version("1.0.0")
	Server("service_new", func() {
		Host("localhost", func() {
			URI("http://0.0.0.0:8081")
			URI("grpc://0.0.0.0:8091")
		})
	})
})
