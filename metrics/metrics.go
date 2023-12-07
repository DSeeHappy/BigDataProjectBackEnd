package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_app_http_requests",
			Help: "Total number of HTTP requests",
		},
	)

	GetJobHttpResponsesCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jobs_app_get_job_http_responses",
			Help: "Total number of HTTP responses for get job API",
		},
		[]string{"status"},
	)

	GetAllJobsTimer = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name: "jobs_app_get_all_jobs_duration",
			Help: "Duration of get all jobs operation",
		},
	)

	GetWeatherHttpResponsesCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jobs_app_get_weather_http_responses",
			Help: "Total number of HTTP responses for get weather API",
		},
		[]string{"status"},
	)
)
