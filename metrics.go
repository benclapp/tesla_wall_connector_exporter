package main

import "github.com/prometheus/client_golang/prometheus"

var (
	uptime = newMetricDesc("uptime_seconds", "Uptime of wallconnector", nil)
)

func newMetricDesc(name string, help string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName("teslawallconnector", "", name), help, labels, nil)
}
