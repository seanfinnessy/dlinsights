package login

import (
	"strings"

	"github.com/yohcop/openid-go"
)

var steamEndpoint = "https://steamcommunity.com/openid"
var simpleDiscoveryCache = openid.NewSimpleDiscoveryCache()
var nonceStore = openid.NewSimpleNonceStore()

// GetRedirectURL returns the Steam OpenID URL to open in the user's browser.
// callbackURL is the loopback address Steam will redirect back to.
func GetRedirectURL(callbackURL string) (string, error) {
	return openid.RedirectURL(steamEndpoint, callbackURL, "")
}

// VerifyCallback verifies a Steam OpenID callback URL and returns the Steam account ID.
func VerifyCallback(fullURL string) (string, error) {
	id, err := openid.Verify(fullURL, simpleDiscoveryCache, nonceStore)
	if err != nil {
		return "", err
	}
	parts := strings.Split(id, "/")
	return parts[len(parts)-1], nil
}
