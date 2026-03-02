package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/seanfinnessy/dlinsights/internal/deadlock"
	"github.com/seanfinnessy/dlinsights/internal/login"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetMatches(steamId string, numMatches int) ([]deadlock.PlayerMatchHistoryEntry, error) {
	return deadlock.GetMatches(steamId, numMatches)
}

func (a *App) GetHeroAssets() ([]deadlock.HeroAssets, error) {
	return deadlock.GetHeroAssets()
}

// SteamLogin starts a loopback HTTP server, opens the Steam login page in the
// system browser, and waits for the OpenID callback. On success it emits
// "steam:authed" with the steamId; on failure it emits "steam:error".
func (a *App) SteamLogin() error {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	callbackURL := fmt.Sprintf("http://localhost:%d/callback", port)

	redirectURL, err := login.GetRedirectURL(callbackURL)
	if err != nil {
		listener.Close()
		return err
	}

	// Mux means multiplex. Matches incoming request URLs to our
	// defined route handlers. Basically a controller.
	mux := http.NewServeMux()
	srv := &http.Server{Handler: mux}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		fullURL := fmt.Sprintf("http://localhost:%d%s", port, r.URL.String())
		steamId, err := login.VerifyCallback(fullURL)
		if err != nil {
			// Emit a steam:error event, with the error data for the frontend
			wailsruntime.EventsEmit(a.ctx, "steam:error", err.Error())
			w.Write([]byte("Authentication failed. You can close this window."))
		} else {
			// Emit a steam:authed event, with the steam id as data for frontend
			wailsruntime.EventsEmit(a.ctx, "steam:authed", steamId)
			w.Write([]byte("Login successful! You can close this window."))
		}
		
		go srv.Shutdown(context.Background())
	})

	go srv.Serve(listener)
	wailsruntime.BrowserOpenURL(a.ctx, redirectURL)

	return nil
}
