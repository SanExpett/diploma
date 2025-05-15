package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// GrpcMetrics представляет собой набор метрик для gRPC с использованием Prometheus
type GrpcMetrics struct {
	// service название отслеживаемого gRPC сервиса
	service string
	// methodCallsTotal счетчик общего количества вызовов методов gRPC
	methodCallsTotal *prometheus.CounterVec
}

// NewGrpcMetrics создает новый экземпляр GrpcMetrics для указанного сервиса
func NewGrpcMetrics(service string) *GrpcMetrics {
	return &GrpcMetrics{
		service: service,
		methodCallsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "method_calls_total",
				Help: "Total amount of grpc methods calls",
			},
			[]string{"service", "method"},
		),
	}
}

// Register регистрирует все метрики в стандартном реестре Prometheus
func (grpcMetrics *GrpcMetrics) Register() {
	prometheus.MustRegister(grpcMetrics.methodCallsTotal)
}

// IncRequestsTotal увеличивает общее количество запросов для указанного метода
func (grpcMetrics *GrpcMetrics) IncRequestsTotal(method string) {
	grpcMetrics.methodCallsTotal.WithLabelValues(grpcMetrics.service, method).Inc()
}
