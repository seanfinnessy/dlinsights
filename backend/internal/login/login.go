package login

import (
	"net/http"
	"strings"

	"github.com/yohcop/openid-go"
)

var steamEndpoint = "https://steamcommunity.com/openid"
var simpleDiscoveryCache = openid.NewSimpleDiscoveryCache()
var nonceStore = openid.NewSimpleNonceStore()

func SteamLogin(w http.ResponseWriter, r *http.Request) {
	redirectURL, err := openid.RedirectURL(
		steamEndpoint,
		"http://localhost:3000/auth/steam/callback",
		"http://localhost:3000/",
	)
	if err != nil {
		http.Error(w, "Failed to redirect", 500)
		return
	}
	
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func SteamCallback(w http.ResponseWriter, r *http.Request) (string, error) {
	fullUrl := "http://localhost:3000" + r.URL.String()
	
	id, err := openid.Verify(
		fullUrl,
		simpleDiscoveryCache, 
		nonceStore,
	)

	if err == nil {
		urlSplice := strings.Split(id, "/")
		lenUrlSplice := len(urlSplice)
		resultingId := urlSplice[lenUrlSplice-1]
		return resultingId, nil

	} else {
		return "", err
	}
}