package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ResponseHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "eta_response_duration_milliseconds",
		Help:    "Duration of ETA responses in milliseconds",
		Buckets: prometheus.LinearBuckets(0.001, 0.01, 100), // Modify buckets as per your needs
	})
)
