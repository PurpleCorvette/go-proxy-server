// Package metrics provides a gRPC unary interceptor for Prometheus metrics.
package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Metrics defines an interface for Prometheus metrics.
type Metrics interface {
	MetricsUnaryInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor
	Handler() http.Handler
}

// PrometheusMetrics implements Metrics interface for Prometheus metrics.
type PrometheusMetrics struct{}

// NewPrometheusMetrics returns a new instance of PrometheusMetrics.
func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{}
}

// requestCount is a Prometheus counter metric for counting gRPC requests.
var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method"},
	)
)

// init registers the requestCount metric with Prometheus.
func init() {
	prometheus.MustRegister(requestCount)
}

// MetricsUnaryInterceptor returns a gRPC unary interceptor for Prometheus metrics.
func (m *PrometheusMetrics) MetricsUnaryInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		requestCount.With(prometheus.Labels{"method": info.FullMethod}).Inc()
		logger.Infof("Metrics updated for method: %s", info.FullMethod)
		return handler(ctx, req)
	}
}

// Handler returns an HTTP handler for serving Prometheus metrics.
func (m *PrometheusMetrics) Handler() http.Handler {
	return promhttp.Handler()
}
