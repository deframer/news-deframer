package design

import . "goa.design/goa/v3/dsl"

var _ = Service("web", func() {
	Description("Browser/web-facing API contract.")
	HTTP(func() {
		Path("/api")
	})

	Error("not_found", String, "Resource not found")
	defineWebMethods()
})

func defineWebMethods() {
	Method("item", func() {
		Description("Fetch a single analyzed item by URL.")
		Payload(func() {
			Attribute("url", String, "Item URL")
			Required("url")
		})
		Result(AnalyzedItem)
		HTTP(func() {
			GET("/item")
			Param("url")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})

	Method("site", func() {
		Description("List analyzed items for a root domain.")
		Payload(func() {
			Attribute("root", String, "Root domain")
			Attribute("max_score", Float64, "Maximum rating to include", func() {
				Default(0)
			})
			Required("root")
		})
		Result(ArrayOf(AnalyzedItem))
		HTTP(func() {
			GET("/site")
			Param("root")
			Param("max_score")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})

	Method("articles", func() {
		Description("List articles for a trend and domain.")
		Payload(func() {
			Attribute("root", String, "Root domain")
			Attribute("term", String, "Trend term")
			Attribute("date", String, "Optional date", func() {
				Format(FormatDate)
			})
			Attribute("days", Int, "Lookback window in days", func() {
				Default(1)
				Minimum(1)
			})
			Attribute("offset", Int, "Result offset", func() {
				Default(0)
				Minimum(0)
			})
			Attribute("limit", Int, "Maximum results", func() {
				Default(20)
				Minimum(1)
			})
			Required("root", "term")
		})
		Result(ArrayOf(AnalyzedArticle))
		HTTP(func() {
			GET("/articles")
			Param("root")
			Param("term")
			Param("date")
			Param("days")
			Param("offset")
			Param("limit")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})

	Method("sentiments", func() {
		Description("Get sentiment scores for a trend.")
		Payload(func() {
			Attribute("root", String, "Root domain")
			Attribute("term", String, "Trend term")
			Attribute("date", String, "Optional date", func() {
				Format(FormatDate)
			})
			Attribute("days", Int, "Lookback window in days", func() {
				Default(1)
				Minimum(1)
			})
			Required("root", "term")
		})
		Result(SentimentItem)
		HTTP(func() {
			GET("/sentiments")
			Param("root")
			Param("term")
			Param("date")
			Param("days")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})

	Method("domains", func() {
		Description("List root domains.")
		Result(ArrayOf(DomainEntry))
		HTTP(func() {
			GET("/domains")
			Response(StatusOK)
		})
	})

	Method("topTrendsByDomain", func() {
		Description("List top trends for a domain.")
		Payload(func() {
			Attribute("domain", String, "Root domain")
			Attribute("lang", String, "Language code")
			Attribute("days", Int, "Lookback window in days", func() {
				Default(1)
				Minimum(1)
			})
			Attribute("date", String, "Optional date", func() {
				Format(FormatDate)
			})
			Required("domain", "lang")
		})
		Result(ArrayOf(TrendMetric))
		HTTP(func() {
			GET("/trends/topbydomain")
			Param("domain")
			Param("lang")
			Param("days")
			Param("date")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})

	Method("contextByDomain", func() {
		Description("List trend context for a domain.")
		Payload(func() {
			Attribute("term", String, "Trend term")
			Attribute("domain", String, "Root domain")
			Attribute("lang", String, "Language code")
			Attribute("days", Int, "Lookback window in days", func() {
				Default(1)
				Minimum(1)
			})
			Attribute("date", String, "Optional date", func() {
				Format(FormatDate)
			})
			Required("term", "domain", "lang")
		})
		Result(ArrayOf(TrendContext))
		HTTP(func() {
			GET("/trends/contextbydomain")
			Param("term")
			Param("domain")
			Param("lang")
			Param("days")
			Param("date")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})

	Method("lifecycleByDomain", func() {
		Description("List trend lifecycle data for a domain.")
		Payload(func() {
			Attribute("term", String, "Trend term")
			Attribute("domain", String, "Root domain")
			Attribute("lang", String, "Language code")
			Attribute("days", Int, "Lookback window in days", func() {
				Default(1)
				Minimum(1)
			})
			Attribute("date", String, "Optional date", func() {
				Format(FormatDate)
			})
			Required("term", "domain", "lang")
		})
		Result(ArrayOf(Lifecycle))
		HTTP(func() {
			GET("/trends/lifecyclebydomain")
			Param("term")
			Param("domain")
			Param("lang")
			Param("days")
			Param("date")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})

	Method("domainComparison", func() {
		Description("Compare two domains for a trend.")
		Payload(func() {
			Attribute("domain_a", String, "First domain")
			Attribute("domain_b", String, "Second domain")
			Attribute("lang", String, "Language code")
			Attribute("days", Int, "Lookback window in days", func() {
				Default(1)
				Minimum(1)
			})
			Attribute("date", String, "Optional date", func() {
				Format(FormatDate)
			})
			Required("domain_a", "domain_b", "lang")
		})
		Result(ArrayOf(DomainComparison))
		HTTP(func() {
			GET("/trends/comparedomains")
			Param("domain_a")
			Param("domain_b")
			Param("lang")
			Param("days")
			Param("date")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})
}
