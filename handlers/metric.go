package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PrometheusMiddleware struct {
	Histogram *prometheus.HistogramVec
	Counter   *prometheus.CounterVec
}

func NewPrometheusMiddleware(registry * prometheus.Registry) *PrometheusMiddleware {

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: "aloha",
		Name:      "request_duration_seconds",
		Help:      "Seconds spent serving HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "aloha",
			Name:      "requests_total",
			Help:      "The total number of HTTP requests.",
		},
		[]string{"status"},
	)

	registry.MustRegister(histogram)
	registry.MustRegister(counter)

	return &PrometheusMiddleware{
		Histogram: histogram,
		Counter:   counter,
	}
}

type metricResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewMetricResponseWriter(w http.ResponseWriter) *metricResponseWriter {
	return &metricResponseWriter{w, http.StatusOK}
}

func (lrw *metricResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (p *PrometheusMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		path := urlToLabel(r.RequestURI)
		mw := NewMetricResponseWriter(w)
		next.ServeHTTP(mw, r)
		var (
			status = strconv.Itoa(mw.statusCode)
			took   = time.Since(begin)
		)
		p.Histogram.WithLabelValues(r.Method, path, status).Observe(took.Seconds())
		p.Counter.WithLabelValues(status).Inc()
	})
}

var invalidChars = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func urlToLabel(path string) string {
	result := invalidChars.ReplaceAllString(path, "_")
	result = strings.ToLower(strings.Trim(result, "_"))
	if result == "" {
		result = "root"
	}
	return result
}