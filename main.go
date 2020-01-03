package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
	"github.com/lukapiske/aloha/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var version = "1.5.0"
var port = ":8090"

func main() {

	log.Print("Starting Aloha!!!")

	prom := handlers.NewPrometheusMiddleware()

	isReady := &atomic.Value{}
	isReady.Store(false)

	go func() {
		log.Printf("Readyz not ready.")
		time.Sleep(10 * time.Second)
		isReady.Store(true)
		log.Printf("Readyz ready.")
	}()

	readyzHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	versionHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(version))
		_, _ = w.Write([]byte("\r\n"))
	})

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Aloha!!!\r\n"))
	})

	healthzHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	badRequestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	forbiddenHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	http.Handle("/healthz", healthzHandler)
	http.Handle("/readyz", readyzHandler)

	http.Handle("/forbidden", prom.Handler(forbiddenHandler))
	http.Handle("/bad", prom.Handler(badRequestHandler))

	http.Handle("/", prom.Handler(rootHandler))
	http.Handle("/version", prom.Handler(versionHandler))
	http.Handle("/metrics", promhttp.HandlerFor(prom.Registry, promhttp.HandlerOpts{}))

	log.Print("Listening on port " + port)

	log.Fatal(http.ListenAndServe(port, nil))

}
