package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	//ResponseHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
	//	Name:    "eta_response_duration_milliseconds",
	//	Help:    "Duration of ETA responses in milliseconds",
	//	Buckets: prometheus.LinearBuckets(0.001, 0.01, 100), // Modify buckets as per your needs
	//})

	ResponseHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_response_time_seconds",
			Help:    "Histogram of response times for gRPC requests",
			Buckets: []float64{0.01, 0.02, 0.03, 0.05, 0.07, 0.08, 0.1, 0.2, 0.5, 1}, // Adjust bucket configuration as needed.
			//Buckets: prometheus.LinearBuckets(0.01, 0.01, 10), // Adjust bucket configuration as needed.
		},
		[]string{"method"},
	)
)
