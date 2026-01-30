package netutil

import (
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// GetRootDomain extracts the effective top-level domain plus one (eTLD+1) from a URL.
// If extraction fails (e.g. for IP addresses or localhost), it falls back to the hostname.
func GetRootDomain(u *url.URL) string {
	if domain, err := publicsuffix.EffectiveTLDPlusOne(u.Hostname()); err == nil {
		return domain
	}
	return u.Hostname()
}

// NormalizeURL removes fragments, common tracking/marketing parameters, and trailing slashes from a URL.
func NormalizeURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	// Remove fragment (e.g., #ref=rss)
	u.Fragment = ""
	u.RawFragment = ""

	// Remove trailing slash from path
	u.Path = strings.TrimSuffix(u.Path, "/")

	// Remove common tracking/marketing parameters (case-insensitive)
	q := u.Query()
	trackingParams := map[string]bool{
		"utm_source": true, "utm_medium": true, "utm_campaign": true,
		"utm_term": true, "utm_content": true, "ref": true,
		"fbclid": true, "gclid": true, "msclkid": true,
		"mc_cid": true, "mc_eid": true, "_hsenc": true, "_hsmi": true,
	}
	for k := range q {
		if trackingParams[strings.ToLower(k)] {
			q.Del(k)
		}
	}

	if len(q) == 0 {
		u.RawQuery = ""
	} else {
		u.RawQuery = q.Encode()
	}

	return u.String()
}
