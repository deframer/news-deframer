package design

import . "goa.design/goa/v3/dsl" //nolint:staticcheck

var HostnameResponse = Type("HostnameResponse", func() {
	Description("Current host name.")
	Attribute("hostname", String, "Machine hostname")
	Required("hostname")
})

var SentimentScores = Type("SentimentScores", func() {
	Description("Sentiment score set.")
	Attribute("valence", Float64, "Valence score")
	Attribute("arousal", Float64, "Arousal score")
	Attribute("dominance", Float64, "Dominance score")
	Attribute("joy", Float64, "Joy score")
	Attribute("anger", Float64, "Anger score")
	Attribute("sadness", Float64, "Sadness score")
	Attribute("fear", Float64, "Fear score")
	Attribute("disgust", Float64, "Disgust score")
})

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

	Method("hostname", func() {
		Description("Return the current host name.")
		Result(HostnameResponse)
		HTTP(func() {
			GET("/hostname")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})
})

var ThinkResult = Type("ThinkResult", func() {
	Description("Analysis result for an article.")
	Attribute("title_original", String, "Original title")
	Attribute("description_original", String, "Original description")
	Attribute("title_corrected", String, "Corrected title")
	Attribute("title_correction_reason", String, "Why the title changed")
	Attribute("description_corrected", String, "Corrected description")
	Attribute("description_correction_reason", String, "Why the description changed")
	Attribute("framing", Float64, "Framing score")
	Attribute("framing_reason", String, "Framing explanation")
	Attribute("clickbait", Float64, "Clickbait score")
	Attribute("clickbait_reason", String, "Clickbait explanation")
	Attribute("persuasive", Float64, "Persuasiveness score")
	Attribute("persuasive_reason", String, "Persuasiveness explanation")
	Attribute("hyper_stimulus", Float64, "Hyper stimulus score")
	Attribute("hyper_stimulus_reason", String, "Hyper stimulus explanation")
	Attribute("speculative", Float64, "Speculative score")
	Attribute("speculative_reason", String, "Speculative explanation")
	Attribute("overall", Float64, "Overall score")
	Attribute("overall_reason", String, "Overall explanation")
})

var MediaThumbnail = Type("MediaThumbnail", func() {
	Description("Media thumbnail metadata.")
	Attribute("url", String, "Thumbnail URL")
	Attribute("height", Int, "Thumbnail height")
	Attribute("width", Int, "Thumbnail width")
	Required("url")
})

var MediaContent = Type("MediaContent", func() {
	Description("Media attachment metadata.")
	Attribute("url", String, "Media URL")
	Attribute("type", String, "MIME type")
	Attribute("medium", String, "Media type")
	Attribute("height", Int, "Media height")
	Attribute("width", Int, "Media width")
	Attribute("title", String, "Media title")
	Attribute("description", String, "Media description")
	Attribute("thumbnail", MediaThumbnail, "Media thumbnail")
	Attribute("credit", String, "Media credit")
	Required("url")
})

var DomainEntry = Type("DomainEntry", func() {
	Description("Root domain entry.")
	Attribute("domain", String, "Root domain")
	Attribute("language", String, "Language code")
	Attribute("portal_url", String, "Portal URL", func() {
		Format(FormatURI)
	})
	Required("domain", "language")
})

var TrendMetric = Type("TrendMetric", func() {
	Description("Trending term metric.")
	Attribute("trend_topic", String, "Trending topic")
	Attribute("frequency", Int64, "Frequency")
	Attribute("utility", Int64, "Utility")
	Attribute("outlier_ratio", Float64, "Outlier ratio")
	Attribute("time_slice", String, "Time slice", func() {
		Format(FormatDateTime)
	})
	Required("trend_topic", "frequency", "utility", "outlier_ratio", "time_slice")
})

var TrendContext = Type("TrendContext", func() {
	Description("Context term for a trend.")
	Attribute("context", String, "Context word")
	Attribute("frequency", Int64, "Frequency")
	Required("context", "frequency")
})

var Lifecycle = Type("Lifecycle", func() {
	Description("Trend lifecycle sample.")
	Attribute("time_slice", String, "Time slice", func() {
		Format(FormatDateTime)
	})
	Attribute("frequency", Int64, "Frequency")
	Attribute("velocity", Int64, "Velocity")
	Required("time_slice", "frequency", "velocity")
})

var DomainComparison = Type("DomainComparison", func() {
	Description("Trend comparison between domains.")
	Attribute("classification", String, "Classification", func() {
		Enum("BLINDSPOT_A", "BLINDSPOT_B", "INTERSECT")
	})
	Attribute("rank_group", Int, "Rank group")
	Attribute("trend_topic", String, "Trending topic")
	Attribute("score_a", Float64, "Score A")
	Attribute("score_b", Float64, "Score B")
	Required("classification", "rank_group", "trend_topic", "score_a", "score_b")
})

var AnalyzedArticle = Type("AnalyzedArticle", func() {
	Description("Article summary for trend views.")
	Attribute("url", String, "Article URL")
	Attribute("title", String, "Article title")
	Attribute("rating", Float64, "Article rating")
	Attribute("authors", ArrayOf(String), "Article authors")
	Attribute("pub_date", String, "Publication date", func() {
		Format(FormatDateTime)
	})
	Required("url", "pub_date")
})

var SentimentItem = Type("SentimentItem", func() {
	Description("Trend sentiment pair.")
	Attribute("sentiments", SentimentScores, "Original sentiments")
	Attribute("sentiments_deframed", SentimentScores, "Deframed sentiments")
})

var AnalyzedItem = Type("AnalyzedItem", func() {
	Description("Full analyzed item.")
	Attribute("hash", String, "Item hash")
	Attribute("url", String, "Item URL")
	Attribute("title_original", String, "Original title")
	Attribute("description_original", String, "Original description")
	Attribute("title_corrected", String, "Corrected title")
	Attribute("title_correction_reason", String, "Why the title changed")
	Attribute("description_corrected", String, "Corrected description")
	Attribute("description_correction_reason", String, "Why the description changed")
	Attribute("framing", Float64, "Framing score")
	Attribute("framing_reason", String, "Framing explanation")
	Attribute("clickbait", Float64, "Clickbait score")
	Attribute("clickbait_reason", String, "Clickbait explanation")
	Attribute("persuasive", Float64, "Persuasiveness score")
	Attribute("persuasive_reason", String, "Persuasiveness explanation")
	Attribute("hyper_stimulus", Float64, "Hyper stimulus score")
	Attribute("hyper_stimulus_reason", String, "Hyper stimulus explanation")
	Attribute("speculative", Float64, "Speculative score")
	Attribute("speculative_reason", String, "Speculative explanation")
	Attribute("overall", Float64, "Overall score")
	Attribute("overall_reason", String, "Overall explanation")
	Attribute("sentiments", SentimentScores, "Original sentiments")
	Attribute("sentiments_deframed", SentimentScores, "Deframed sentiments")
	Attribute("media", MediaContent, "Media content")
	Attribute("rating", Float64, "Think rating")
	Attribute("authors", ArrayOf(String), "Authors")
	Attribute("pubDate", String, "Publication date", func() {
		Format(FormatDateTime)
	})
	Required("hash", "url", "rating", "pubDate", "authors")
})

var AnalyzedSiteItem = Type("AnalyzedSiteItem", func() {
	Description("Analyzed item for site listings.")
	Attribute("hash", String, "Item hash")
	Attribute("url", String, "Item URL")
	Attribute("title_original", String, "Original title")
	Attribute("description_original", String, "Original description")
	Attribute("title_corrected", String, "Corrected title")
	Attribute("title_correction_reason", String, "Why the title changed")
	Attribute("description_corrected", String, "Corrected description")
	Attribute("description_correction_reason", String, "Why the description changed")
	Attribute("framing", Float64, "Framing score")
	Attribute("framing_reason", String, "Framing explanation")
	Attribute("clickbait", Float64, "Clickbait score")
	Attribute("clickbait_reason", String, "Clickbait explanation")
	Attribute("persuasive", Float64, "Persuasiveness score")
	Attribute("persuasive_reason", String, "Persuasiveness explanation")
	Attribute("hyper_stimulus", Float64, "Hyper stimulus score")
	Attribute("hyper_stimulus_reason", String, "Hyper stimulus explanation")
	Attribute("speculative", Float64, "Speculative score")
	Attribute("speculative_reason", String, "Speculative explanation")
	Attribute("overall", Float64, "Overall score")
	Attribute("overall_reason", String, "Overall explanation")
	Attribute("media", MediaContent, "Media content")
	Attribute("rating", Float64, "Think rating")
	Attribute("authors", ArrayOf(String), "Authors")
	Attribute("pubDate", String, "Publication date", func() {
		Format(FormatDateTime)
	})
	Required("hash", "url", "rating", "pubDate", "authors")
})
