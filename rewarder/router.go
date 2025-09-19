package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Router initializes and returns a new mux.Router with defined routes and middleware.
func Router() *mux.Router {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	// apply middleware
	router.Use(corsMiddleware)

	// register route
	router.HandleFunc("/upload", upload).Methods("POST")
	//router.HandleFunc("/download/{companyName}", download).Methods("GET")

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		// optional: send 404 or some default response for other unmatched routes
		http.NotFound(w, r)
	})

	return router
}

// corsMiddleware adds CORS headers to the response
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
