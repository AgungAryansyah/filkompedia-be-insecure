package prometheus

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Info            *prometheus.GaugeVec
	RequestTotal    *prometheus.CounterVec
	Duration        *prometheus.HistogramVec
	DurationSummary prometheus.Summary
	ErrorCount      *prometheus.CounterVec
}

func PrometheusNewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		Info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "filkompedia_be",
			Name:      "info",
			Help:      "",
		},
			[]string{"version"}),
		RequestTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "filkompedia_be",
			Name:      "number_total_request",
			Help:      "Number of total request",
		}, []string{"response_code", "method"}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "filkompedia_be",
			Name:      "request_duration_seconds",
			Buckets:   []float64{.1, .15, .2, .25, .3},
			Help:      "Request duration in seconds",
		}, []string{"status", "method", "route"}),
		DurationSummary: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  "filkompedia_be",
			Name:       "request_duration_summary_seconds",
			Help:       "Request duration in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		ErrorCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "filkompedia_be",
			Name:      "error_count",
			Help:      "Number of error responses",
		}, []string{"method", "response_code"}),
	}

	reg.MustRegister(m.Info, m.RequestTotal, m.Duration, m.DurationSummary)
	return m
}

func Start() Metrics {
	promReg := prometheus.NewRegistry()
	promReg.MustRegister(collectors.NewGoCollector())

	promMetrics := PrometheusNewMetrics(promReg)
	promMetrics.Info.WithLabelValues(runtime.Version()).Set(1)

	pMux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{})
	pMux.Handle("/metrics", promHandler)

	go func() {
		port := os.Getenv("PROMETHEUS_EXPORTER_PORT")
		log.Printf("Starting prometheus exporter on port %s", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), pMux); err != nil {
			log.Fatal("Error when running prometheus exporter, error: " + err.Error())
		}
	}()

	return *promMetrics
}
