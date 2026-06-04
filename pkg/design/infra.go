package design

import . "goa.design/goa/v3/dsl" //nolint:staticcheck

var _ = Service("infra", func() {
	Description("Infrastructure and health endpoints.")
	HTTP(func() {
		Path("")
	})

	Error("not_found", String, "Resource not found")

	Method("ping", func() {
		Description("Health check endpoint.")
		Result(String)
		HTTP(func() {
			GET("/ping")
			Response(StatusOK, func() {
				ContentType("text/plain")
			})
		})
	})
})
