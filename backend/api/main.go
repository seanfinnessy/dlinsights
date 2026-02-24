package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/seanfinnessy/dlinsights/internal/deadlock"
	"github.com/seanfinnessy/dlinsights/internal/login"
)

// using this to hold id for now
var tempHoldId string
func main() {
	// Create a new router and set up logger
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/auth/steam/login", func(w http.ResponseWriter, r *http.Request) {
		login.SteamLogin(w, r)
	})

	r.Get("/auth/steam/callback", func(w http.ResponseWriter, r *http.Request) {
		if id, err := login.SteamCallback(w, r); err == nil {
			tempHoldId = id
			fmt.Println("Saved id locally: " + tempHoldId)
			http.Redirect(w, r, "http://localhost:5173?steamId="+id, http.StatusFound)
		} else {
			http.Redirect(w, r, "http://localhost:5173?error=auth_failed", http.StatusFound)
		}
	})

	// Get history of matches for given user id
	r.Get("/get-matches/{userId}", func(w http.ResponseWriter, r *http.Request) {
		// Grab URL and Query params
		userId := chi.URLParam(r, "userId")
		numMatchesStr := r.URL.Query().Get("amount")
		if numMatchesStr == "" {
			numMatchesStr = "10" // fallback, URL.Query().Get() returns empty string if not found
		}
		numMatches, err := strconv.Atoi(numMatchesStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error: "+err.Error())
		}

		matches, err := deadlock.GetMatches(userId, numMatches)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error: "+err.Error())
			return
		}

		// Set headers
		w.Header().Set("Content-Type", "application/json")
		// Stream JSON directly into response writer
		json.NewEncoder(w).Encode(matches)

	})

	r.Get("/get-hero-assets", func(w http.ResponseWriter, r *http.Request) {
		heroAssets, err := deadlock.GetHeroAssets()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(heroAssets)
	})

	// Setup CORS handler
	handler := cors.Default().Handler(r)
	
	fmt.Println("Serving...")
	http.ListenAndServe(":3000", handler)
}