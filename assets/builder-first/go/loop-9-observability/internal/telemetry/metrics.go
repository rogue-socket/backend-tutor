// RED metrics for HTTP traffic, plus DB query histogram for tier-2 visibility.
//
// CARDINALITY DISCIPLINE:
//   - Labels: low-cardinality only. method, route (not path!), status.
//   - NEVER: user_id, request_id, IP, raw URL path with IDs.
//   - High-cardinality data goes in logs and traces, not metrics.

package telemetry

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests, partitioned by method, route, status.",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "route", "status"},
	)

	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "DB query duration in seconds, by query name.",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5},
		},
		[]string{"query"}, // query name, NOT the SQL itself (high-card)
	)

	wsClientsConnected = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ws_clients_connected",
			Help: "Currently connected WebSocket clients.",
		},
	)

	queueDepth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "queue_depth",
			Help: "Pending messages per queue.",
		},
		[]string{"stream"},
	)
)

// Middleware wraps handlers, recording request count and duration with the
// matched route (not the raw path).
//
// Pass route as a literal — "/links/{id}", not r.URL.Path. The middleware
// can't infer the route from the URL because that's the high-cardinality
// trap we're avoiding.
func Middleware(route string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)
		elapsed := time.Since(start).Seconds()
		labels := prometheus.Labels{
			"method": r.Method,
			"route":  route,
			"status": strconv.Itoa(ww.status),
		}
		httpRequestsTotal.With(labels).Inc()
		httpRequestDuration.With(labels).Observe(elapsed)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// ObserveDBQuery is called from the data-access layer.
func ObserveDBQuery(name string, d time.Duration) {
	dbQueryDuration.WithLabelValues(name).Observe(d.Seconds())
}

// SetWSClients reflects the current count from the Hub.
func SetWSClients(n int) { wsClientsConnected.Set(float64(n)) }

// SetQueueDepth from the queue worker's monitoring tick.
func SetQueueDepth(stream string, n int) { queueDepth.WithLabelValues(stream).Set(float64(n)) }
