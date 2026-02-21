package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/seanfinnessy/dlinsights/internal/login"
	"github.com/seanfinnessy/dlinsights/internal/deadlock"
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
		}
	})
	r.Get("/get-matches/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		fmt.Println("Requested ID:", id)

		err := deadlock.GetMatches(id)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error: "+err.Error())
			return
		}

		fmt.Fprint(w, "Success for ID "+id)
	})
	fmt.Println("Starting server...")
	http.ListenAndServe(":3000", r)
}