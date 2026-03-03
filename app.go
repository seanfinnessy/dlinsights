package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/seanfinnessy/dlinsights/internal/deadlock"
	"github.com/seanfinnessy/dlinsights/internal/login"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func addWatcherRecursive(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return watcher.Add(path)
		}

		return nil
	})
}

func getDLchanges() {
	rootPath := "C:\\Program Files (x86)\\Steam\\steamapps\\common\\Deadlock"

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Watch existing folders
	err = addWatcherRecursive(watcher, rootPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Watching all folders under:", rootPath)

	for {
		select {

		case event := <-watcher.Events:

			fmt.Println("Event:", event)

			// New directory → watch it
			if event.Op&fsnotify.Create == fsnotify.Create {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					watcher.Add(event.Name)
					fmt.Println("New directory watched:", event.Name)
				}
			}

			// File written → read contents
			if event.Op&fsnotify.Write == fsnotify.Write {

				if info, err := os.Stat(event.Name); err == nil && !info.IsDir() {
					if event.Name.includes("reconnect.dat") { // TODO: "C:\\Program Files (x86)\\Steam\\steamapps\\common\\Deadlock\\game\\citadel\\reconnect.dat" is created and written at game start, and removed at end
						fmt.Println("Found game")
					}
					data, err := os.ReadFile(event.Name)
					if err == nil {
						fmt.Println("File updated:", event.Name)
						fmt.Println(string(data))
					}
				}
			}

		case err := <-watcher.Errors:
			fmt.Println("Watcher error:", err)
		}
	}
}

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	getDLchanges()
}

func (a *App) GetMatches(steamId string, numMatches int) ([]deadlock.PlayerMatchHistoryEntry, error) {
	return deadlock.GetMatches(steamId, numMatches)
}

func (a *App) GetMatchInfo(matchId string) (deadlock.MatchInfoResponse, error) {
	return deadlock.GetMatchInfo(matchId)
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
