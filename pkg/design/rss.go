package design

import . "goa.design/goa/v3/dsl" //nolint:staticcheck

var RSSPayload = Type("RSSPayload", func() {
	Description("RSS proxy request.")
	Attribute("url", String, "Feed URL")
	Attribute("lang", String, "Language code")
	Attribute("max_score", Float64, "Maximum rating to include")
	Attribute("embedded", Boolean, "Embed content in the feed")
	Required("url")
})

var _ = Service("rss", func() {
	Description("RSS proxy endpoint.")
	HTTP(func() {
		Path("")
	})

	Error("bad_request", String, "Bad request")

	Method("feed", func() {
		Description("Return a proxied RSS feed.")
		Payload(RSSPayload)
		Result(String)
		HTTP(func() {
			GET("/rss")
			Param("url")
			Param("lang")
			Param("max_score")
			Param("embedded")
			Response(StatusOK, func() {
				ContentType("text/plain")
			})
			Response("bad_request", StatusBadRequest)
		})
	})
})
