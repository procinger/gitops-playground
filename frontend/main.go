package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	backendHost := os.Getenv("BACKEND_URL")
	if backendHost == "" {
		backendHost = "http://backend:8080"
	}

	backendURL, err := url.Parse(backendHost)
	if err != nil {
		log.Fatalf("invalid BACKEND_URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	mux := http.NewServeMux()

	mux.Handle("/api/", proxy)

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "ok")
	})

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Frontend reachable at %s (Backend: %s)", addr, backendHost)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server Error: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
