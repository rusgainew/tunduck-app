package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler предоставляет HTTP handler для Prometheus metrics
func MetricsHandler() any {
	return promhttp.Handler()
}

// Metrics содержит все Prometheus метрики приложения
type Metrics struct {
	// HTTP метрики
	HTTPRequestsTotal   prometheus.Counter
	HTTPRequestDuration prometheus.Histogram
	HTTPRequestSize     prometheus.Histogram
	HTTPResponseSize    prometheus.Histogram

	// Cache метрики
	CacheHitsTotal         prometheus.Counter
	CacheMissesTotal       prometheus.Counter
	CacheEvictionsTotal    prometheus.Counter
	CacheItemsTotal        prometheus.Gauge
	CacheOperationDuration prometheus.Histogram

	// Database метрики
	DBQueryDuration     prometheus.Histogram
	DBErrorsTotal       prometheus.Counter
	DBConnectionsActive prometheus.Gauge

	// Auth метрики
	LoginAttemptsTotal  prometheus.Counter
	LogoutAttemptsTotal prometheus.Counter
	TokensRevokedTotal  prometheus.Counter
	AuthErrorsTotal     prometheus.Counter

	// Business метрики
	UsersRegisteredTotal      prometheus.Counter
	DocumentsCreatedTotal     prometheus.Counter
	OrganizationsCreatedTotal prometheus.Counter
}

// NewMetrics создает и инициализирует все Prometheus метрики
func NewMetrics() *Metrics {
	return &Metrics{
		// HTTP метрики
		HTTPRequestsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		}),
		HTTPRequestDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		}),
		HTTPRequestSize: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		}),
		HTTPResponseSize: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		}),

		// Cache метрики
		CacheHitsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		}),
		CacheMissesTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		}),
		CacheEvictionsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "cache_evictions_total",
			Help: "Total number of cache evictions",
		}),
		CacheItemsTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cache_items_total",
			Help: "Total number of items in cache",
		}),
		CacheOperationDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "cache_operation_duration_seconds",
			Help:    "Cache operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		}),

		// Database метрики
		DBQueryDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		}),
		DBErrorsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "db_errors_total",
			Help: "Total number of database errors",
		}),
		DBConnectionsActive: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		}),

		// Auth метрики
		LoginAttemptsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "auth_login_attempts_total",
			Help: "Total number of login attempts",
		}),
		LogoutAttemptsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "auth_logout_attempts_total",
			Help: "Total number of logout attempts",
		}),
		TokensRevokedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "auth_tokens_revoked_total",
			Help: "Total number of revoked tokens",
		}),
		AuthErrorsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "auth_errors_total",
			Help: "Total number of authentication errors",
		}),

		// Business метрики
		UsersRegisteredTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "users_registered_total",
			Help: "Total number of registered users",
		}),
		DocumentsCreatedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "documents_created_total",
			Help: "Total number of created documents",
		}),
		OrganizationsCreatedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "organizations_created_total",
			Help: "Total number of created organizations",
		}),
	}
}
