package main

import "github.com/prometheus/client_golang/prometheus"

var (
	up = newMetricDesc("up", "Status of Tesla Wall Connector API calls", []string{"collector"})
	// Lifetime
	chargeStarts = newMetricDesc("charge_starts_total", "Number of charges started (TBC)", nil)
	chargingTime = newMetricDesc("charging_duration_seconds", "Total time spent charging (TBC)", nil)
	energyWh     = newMetricDesc("delivered_energy_watt_hours_total", "Total energy delivered via Wall connector (TBC)", nil)
	uptime       = newMetricDesc("uptime_seconds", "Uptime of wallconnector", nil)

	// Version
	info = newMetricDesc("info", "Version and general info about the Tesla Wall Connector",
		[]string{"firmware_version", "part_number", "serial_number"})

	// Vitals
	vehicleConnected = newMetricDesc("vehicle_connected", "Whether or not a vehicle is connected", nil)
)

func newMetricDesc(name string, help string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName("teslawallconnector", "", name), help, labels, nil)
}
