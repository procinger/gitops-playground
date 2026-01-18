package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

const DefaultBackendHost = "http://backend:8080"

func main() {
	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion = "none"
	}

	backendHost := os.Getenv("BACKEND_URL")
	if backendHost == "" {
		backendHost = DefaultBackendHost
	}

	proxy, err := NewProxy(backendHost)
	if err != nil {
		log.Fatalf("failed to create proxy: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", proxy)

	static := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", static))

	tmpl := template.Must(template.ParseFiles("./template/index.html"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"AppVersion": appVersion,
		}
		tmpl.Execute(w, data)
	})

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

	log.Printf("Frontend reachable at %s (Backend at: %s)", addr, backendHost)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server Error: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %s %s", r.Method, r.UserAgent(), r.URL.Path, time.Since(start), r.Header)
	})
}

func NewProxy(host string) (*httputil.ReverseProxy, error) {
	backendURL, err := url.Parse(host)
	if err != nil {
		log.Fatalf("invalid BACKEND_URL: %v", err)
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Scheme = backendURL.Scheme
		req.URL.Host = backendURL.Host
		req.Host = backendURL.Host
	}

	return proxy, nil
}
