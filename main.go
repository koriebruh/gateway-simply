package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/koriebruh/gateway-simply/handlers"
	"github.com/koriebruh/gateway-simply/test"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"path"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	// Register metrik ke Prometheus
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(requestDuration)
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome Home")
}

func main() {
	r := mux.NewRouter()

	// Middleware untuk menghitung metrik
	r.Use(promMiddleware)

	r.PathPrefix("/products").HandlerFunc(handlers.ProxyRequest)
	r.PathPrefix("/orders").HandlerFunc(handlers.ProxyRequest)

	r.HandleFunc("/", Home)

	// Endpoint untuk Prometheus scrapes
	r.Handle("/metrics", promhttp.Handler())

	log.Print("RUNNING ON PORT 8080")
	go test.DummyClient("products", "8081")
	go test.DummyClient("orders", "8082")

	log.Fatal(http.ListenAndServe(":8080", r))
}

// Middleware untuk mengumpulkan metrik
func promMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(requestDuration.WithLabelValues(path))
		totalRequests.WithLabelValues(path).Inc()

		next.ServeHTTP(w, r)

		timer.ObserveDuration()
	})
}
