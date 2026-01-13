package netutil

import (
	"net/url"

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
