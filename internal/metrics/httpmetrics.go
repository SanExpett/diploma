package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// HttpMetrics структура для хранения метрик HTTP-запросов
type HttpMetrics struct {
	requestsTotal   *prometheus.CounterVec   // Счетчик общего количества HTTP-запросов
	requestDuration *prometheus.HistogramVec // Гистограмма длительности HTTP-запросов
}

// NewHttpMetrics создает новый экземпляр структуры HttpMetrics с инициализированными метриками
func NewHttpMetrics() *HttpMetrics {
	return &HttpMetrics{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total amount of http requests with status codes",
			},
			[]string{"endpoint", "method", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_requests_duration",
				Help:    "Duration of http requests with status codes",
				Buckets: []float64{0.001, 0.01, 0.1, 1, 10, 100, 1000, 10000},
			},
			[]string{"endpoint", "method"},
		),
	}
}

// Register регистрирует метрики в системе Prometheus
func (httpMetrics *HttpMetrics) Register() {
	prometheus.MustRegister(httpMetrics.requestsTotal)
	prometheus.MustRegister(httpMetrics.requestDuration)
}

// IncRequestsTotal увеличивает счетчик общего количества запросов для указанного эндпоинта, метода и статуса
func (httpMetrics *HttpMetrics) IncRequestsTotal(endpoint, method string, status int) {
	httpMetrics.requestsTotal.WithLabelValues(endpoint, method, fmt.Sprintf("%d", status)).Inc()
}

// IncRequestDuration добавляет значение длительности запроса в гистограмму для указанного эндпоинта и метода
func (httpMetrics *HttpMetrics) IncRequestDuration(endpoint, method string, duration float64) {
	httpMetrics.requestDuration.WithLabelValues(endpoint, method).Observe(duration)
}
