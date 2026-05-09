// Package design goa service DSL
package design

import (
	. "goa.design/goa/v3/dsl" // nolint:staticcheck
)

var _ = Service("openapi", func() {
	Description("This service provides openapi functions.")

	Method("openApi", func() {
		Payload(func() {
			Field(1, "forwarded_host", String, "X-Forwarded-Host")
			Field(2, "forwarded_proto", String, "X-Forwarded-Proto")
			Field(3, "tag", String, "Tag")
		})

		Result(String)

		Error("invalid_openapi_file", ErrorResult, "Can't parse the openapi file")

		HTTP(func() {
			GET("/openapi/{tag}")
			Headers(func() {
				Header("forwarded_host:X-Forwarded-Host", String, "X-Forwarded-Host", func() {
				})
				Header("forwarded_proto:X-Forwarded-Proto", String, "X-Forwarded-Proto", func() {
				})
			})

			Response(StatusOK, func() {
				ContentType("text/plain")
			})

			Response("invalid_openapi_file", StatusInternalServerError)
		})
	})

	Method("openApiOnlyFor", func() {
		Payload(func() {
			Field(1, "forwarded_host", String, "X-Forwarded-Host")
			Field(2, "forwarded_proto", String, "X-Forwarded-Proto")
			Field(3, "tag", String, "Tag")
			Field(4, "paths", ArrayOf(String), "Only for this paths")
		})

		Result(String)

		Error("invalid_openapi_file", ErrorResult, "Can't parse the openapi file")

		HTTP(func() {
			POST("/openapi")

			Headers(func() {
				Header("forwarded_host:X-Forwarded-Host", String, "X-Forwarded-Host", func() {
				})
				Header("forwarded_proto:X-Forwarded-Proto", String, "X-Forwarded-Proto", func() {
				})
			})

			Response(StatusOK, func() {
				ContentType("text/plain")
			})

			Response("invalid_openapi_file", StatusInternalServerError)
		})
	})
})
