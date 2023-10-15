package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress = flag.String("web.listen-address", ":9859",
		"Address to listen on for HTTP requests.")
	metricsPath = flag.String("web.metrics-path", "/metrics",
		"Path to expose metrics on.")
	twcAddress = flag.String("twc.address", "",
		"[REQUIRED] The address of the Tesla Wall Connector.")
	twcScrapeTimeout = flag.Duration("twc.scrape-timeout", time.Second,
		"The timeout for the scrape request.")

	// Provided at build time
	builtBy, commit, date, version string
)

func main() {
	flag.Parse()
	if len(*twcAddress) == 0 {
		log.Fatal("Address for the Tesla Wall Connector is required.")
	}

	slog.Info("Tesla Wall Connector Exporter", "builtBy", builtBy, "commit", commit, "date", date, "version", version)
	slog.Info("Configuration: ",
		"web.listen-address", *listenAddress,
		"web.metrics-path", *metricsPath,
		"twc.address", *twcAddress,
		"twc.scrape-timeout", *twcScrapeTimeout,
	)

	var build_info = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "teslawallconnector_build_info",
		Help: "Build info about the tessla_wall_connector_exporter",
		ConstLabels: prometheus.Labels{
			"builtBy": builtBy,
			"commit":  commit,
			"date":    date,
			"version": version,
		}})
	prometheus.MustRegister(build_info, NewExporter())
	build_info.Set(1)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Tesla Wall Connector Exporter</title></head>
			<body>
			<h1>Tesla Wall Connector Exporter</h1>
			<p><a href=` + *metricsPath + `>Metrics</a></p>
			</body>
			</html>`))
	})

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
