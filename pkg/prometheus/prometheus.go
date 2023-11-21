package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ResponseHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "eta_response_duration_seconds",
		Help:    "Duration of ETA responses in seconds",
		Buckets: prometheus.LinearBuckets(0.001, 0.01, 100), // Modify buckets as per your needs
	})
)
