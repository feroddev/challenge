package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duração das requisições HTTP em segundos",
		Buckets: []float64{0.1, 0.3, 0.5, 0.7, 1, 2, 5, 10},
	}, []string{"method", "endpoint", "status"})

	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total de requisições HTTP",
	}, []string{"method", "endpoint", "status"})

	DatabaseOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "db_operation_duration_seconds",
		Help: "Duração das operações de banco de dados em segundos",
		Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2},
	}, []string{"operation", "entity"})

	NatsPublishDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "nats_publish_duration_seconds",
		Help: "Duração das operações de publicação no NATS em segundos",
		Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1},
	}, []string{"topic"})

	CacheOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "cache_operation_duration_seconds",
		Help: "Duração das operações de cache em segundos",
		Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1},
	}, []string{"operation"})
)

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func MeasureRequestDuration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start).Seconds()
		statusCode := rw.statusCode
		
		RequestDuration.WithLabelValues(r.Method, r.URL.Path, string(rune(statusCode))).Observe(duration)
		RequestsTotal.WithLabelValues(r.Method, r.URL.Path, string(rune(statusCode))).Inc()
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
