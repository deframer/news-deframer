package design

import . "goa.design/goa/v3/dsl"

var _ = Service("mobile", func() {
	Description("Mobile-facing API contract.")
	HTTP(func() {
		Path("/mobile/api")
	})

	Error("not_found", String, "Resource not found")
	defineWebMethods()
})
