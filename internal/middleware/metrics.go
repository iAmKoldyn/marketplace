package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of HTTP requests",
		},
		[]string{"method", "route", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
}

// Instrument handler to count requests
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		route := middleware.GetRouteContext(r).RoutePattern()
		httpRequestsTotal.WithLabelValues(r.Method, route, http.StatusText(ww.Status())).Inc()
	})
}

// Expose /metrics endpoint in your router
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
