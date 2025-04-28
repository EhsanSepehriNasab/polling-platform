package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Basic HTTP metrics
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for handlers in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Cache metrics
	CacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// DB query metrics
	DBQueryDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Histogram of DB query durations in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	// Specific handler metrics (Poll Feed)
	PollFeedDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "poll_feed_duration_seconds",
			Help:    "Duration of PollFeedHandler execution in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	// Register all metrics with Prometheus
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(CacheHits)
	prometheus.MustRegister(CacheMisses)
	prometheus.MustRegister(DBQueryDuration)
	prometheus.MustRegister(PollFeedDuration)
}

// MetricsMiddleware tracks request count and duration for every HTTP request
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture the status code
		rw := &statusResponseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		// Update global metrics
		HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rw.status)).Inc()
		HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

// Helper to capture status code properly
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *statusResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
