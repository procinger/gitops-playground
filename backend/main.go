package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type TimeResponse struct {
	ISO  string `json:"iso"`
	Unix int64  `json:"unix"`
	Zone string `json:"zone"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		zone, _ := now.Zone()
		resp := TimeResponse{
			ISO:  now.UTC().Format(time.RFC3339Nano),
			Unix: now.Unix(),
			Zone: zone,
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(resp)
	})

	cors := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			h.ServeHTTP(w, r)
		})
	}

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           cors(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Service is running at %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}
