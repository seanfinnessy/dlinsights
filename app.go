package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/seanfinnessy/dlinsights/internal/deadlock"
	"github.com/seanfinnessy/dlinsights/internal/login"
	"github.com/shirou/gopsutil/v4/process"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func isGameRunning() bool {
	processes, _ := process.Processes()
	for _, p := range processes {
		name, _ := p.Name()
		if strings.Contains(strings.ToLower(name), "deadlock") {
			return true
		}
	}
	
	return false
}

func checkForActiveGame(a *App, stop chan struct{}) {
	path := "C:\\Program Files (x86)\\Steam\\steamapps\\common\\Deadlock\\game\\citadel\\"

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Attach our watcher to citadel path to look for Writing to reconnect.dat
	// Created and written to when game is found
	// Removed when leaving game
	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
			case event := <-watcher.Events:

				fmt.Println("Event:", event)

				// File written → read contents
				// event.Op is a bitmask, fsnotify defines file op types as bit flags.
				// so here we are just comparing the bits of event.Op to the bits of Write. If it contains the write bits, it is a write event
				// "is the write bit turned on inside of event.Op"
				if event.Op & fsnotify.Write == fsnotify.Write {
					// verify the written to file isnt a dir..
					if info, err := os.Stat(event.Name); err == nil && !info.IsDir() {
						if strings.Contains(event.Name, "reconnect.dat") { 
							fmt.Println("Found game!")
							a.status.SetInProgressMatch(true)
						}
					}
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					if info, err := os.Stat(event.Name); err == nil && !info.IsDir() {
						if strings.Contains(event.Name, "reconnect.dat") { 
							fmt.Println("Left game!")
							a.status.SetInProgressMatch(false)
						}
					}
				}

			case err := <-watcher.Errors:
				fmt.Println("Watcher error:", err)
				a.status.SetInProgressMatch(false)

			case <-stop:
				fmt.Println("Game closed, stopping file watcher")
				a.status.SetInProgressMatch(false)
				return
			}
	}
}

type App struct {
	ctx context.Context
	status GameStatus
}

type GameStatus struct {
	GameRunning bool
	InProgressMatch bool
	mu sync.RWMutex
}

func NewApp() *App {
	return &App{}
}


func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	go func() {
		// create ticker
		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()

		// create a channel on initial start..
		stopWatcher := make(chan struct{})
		var watcherIsRunning bool

		// monitor game state
		for range(ticker.C) {
			gameRunning := isGameRunning()
			if gameRunning && !watcherIsRunning {
				fmt.Println("Game is now open, starting the file watcher...")
				// recreate the channel if game is restarted
				stopWatcher := make(chan struct{})
				a.status.SetGameRunning(true)
				watcherIsRunning = true
				go checkForActiveGame(a, stopWatcher)
				
			}

			if !gameRunning && watcherIsRunning {
				fmt.Println("Game is now closed, stop the file watcher...")
				a.status.SetGameRunning(false)
				// close channel  when game is closed
				close(stopWatcher)
				watcherIsRunning = false
			}
		}
	}()

}

func (s *GameStatus) SetGameRunning(running bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.GameRunning = running
}

func (s *GameStatus) SetInProgressMatch(inMatch bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.InProgressMatch = inMatch
}

func (s *GameStatus) Get() (bool, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.GameRunning, s.InProgressMatch
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
