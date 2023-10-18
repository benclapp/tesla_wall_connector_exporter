package main

import "github.com/prometheus/client_golang/prometheus"

var (
	up             = newMetricDesc("up", "Status of Tesla Wall Connector API calls", []string{"collector"})
	scrapeDuration = newMetricDesc("scrape_duration_seconds", "Duration of scrapes per collector", []string{"collector"})
	// Lifetime
	chargeStarts          = newMetricDesc("lifetime_charge_starts_total", "Number of charges started (TBC)", nil)
	chargingTime          = newMetricDesc("lifetime_charging_duration_seconds", "Total time spent charging (TBC)", nil)
	energyWh              = newMetricDesc("lifetime_energy_watt_hours_total", "Total energy delivered via Wall connector (TBC)", nil)
	uptime                = newMetricDesc("uptime_seconds", "Uptime of wallconnector", nil)
	contactorCycles       = newMetricDesc("lifetime_contactor_cycles_total", "Number of contactor cycles", nil)
	contactorCyclesLoaded = newMetricDesc("lifetime_contactor_cycles_loaded_total", "Number of contactor cycles while loaded", nil)
	lifetimeAlerts        = newMetricDesc("lifetime_alerts_total", "Number of lifetime alerts", nil)
	thermalFoldbacks      = newMetricDesc("lifetime_thermal_foldbacks_total", "Number of thermal foldbacks", nil)

	// Version
	info = newMetricDesc("info", "Version and general info about the Tesla Wall Connector",
		[]string{"firmware_version", "part_number", "serial_number"})

	// Vitals
	alertsCount        = newMetricDesc("current_alerts", "How many current alerts are there", nil)
	gridHz             = newMetricDesc("grid_hertz", "Current Grid Frequency", nil)
	gridV              = newMetricDesc("grid_voltage", "Current Grid Voltage", nil)
	currentAA          = newMetricDesc("currentA_amps", "Current A amps", nil)
	currentBA          = newMetricDesc("currentB_amps", "Current B amps", nil)
	currentCA          = newMetricDesc("currentC_amps", "Current C amps", nil)
	currentNA          = newMetricDesc("currentN_amps", "Current N amps", nil)
	voltageAV          = newMetricDesc("voltageA_volts", "Voltage A volts", nil)
	voltageBV          = newMetricDesc("voltageB_volts", "Voltage B volts", nil)
	voltageCV          = newMetricDesc("voltageC_volts", "Voltage C volts", nil)
	handleTemp         = newMetricDesc("handle_temperature_celsius", "Current temperature of the handle", nil)
	mcuTemp            = newMetricDesc("mcu_temperature_celsius", "Current temperature of the main control unit", nil)
	pcbaTemp           = newMetricDesc("pcba_temperature_celsius", "Current temperature of PCBA", nil)
	sessionEnergyWh    = newMetricDesc("session_energy_watt_hours_total", "Energy delivered in the current charge session (TBC)", nil)
	sessionSeconds     = newMetricDesc("session_duration_seconds", "Charge session duration in seconds (TBC)", nil)
	vehicleConnected   = newMetricDesc("vehicle_connected", "Whether or not a vehicle is connected", nil)
	vehicleCurrentAmps = newMetricDesc("vehicle_current_amps", "Amps drawn by the connected vehicle (TBC)", nil)
)

func newMetricDesc(name string, help string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName("teslawallconnector", "", name), help, labels, nil)
}
