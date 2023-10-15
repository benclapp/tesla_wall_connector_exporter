package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const apiLifetime = "/api/1/lifetime"
const apiVersion = "/api/1/version"
const apiVitals = "/api/1/vitals"

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	// Lifetime
	ch <- chargeStarts
	ch <- chargingTime
	ch <- energyWh
	ch <- uptime
	// Version
	ch <- info
	// Vitals
	ch <- alertsCount
	ch <- gridHz
	ch <- gridV
	ch <- handleTemp
	ch <- mcuTemp
	ch <- pcbaTemp
	ch <- sessionEnergyWh
	ch <- sessionSeconds
	ch <- vehicleConnected
	ch <- vehicleCurrentAmps
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	client := http.Client{
		Timeout: *twcScrapeTimeout,
	}

	// Scrape the Wall Connector, unmarshal to structs. Update up metric with status
	lt, ltT, err := scrapeLifetime(client)
	if err != nil {
		slog.Error("Failed collection", "err", err, "path", apiLifetime)
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0, apiLifetime)
	} else {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1, apiLifetime)

		ch <- prometheus.MustNewConstMetric(chargeStarts, prometheus.GaugeValue, float64(lt.ChargeStarts))
		ch <- prometheus.MustNewConstMetric(chargingTime, prometheus.GaugeValue, float64(lt.ChargingTimeS))
		ch <- prometheus.MustNewConstMetric(energyWh, prometheus.GaugeValue, float64(lt.EnergyWh))
		ch <- prometheus.MustNewConstMetric(uptime, prometheus.GaugeValue, float64(lt.UptimeS))
	}
	ch <- prometheus.MustNewConstMetric(scrapeDuration, prometheus.GaugeValue, ltT, apiLifetime)

	version, verT, err := scrapeVersion(client)
	if err != nil {
		slog.Error("Failed collection", "err", err, "path", apiVersion)
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0, apiVersion)
	} else {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1, apiVersion)

		ch <- prometheus.MustNewConstMetric(info, prometheus.GaugeValue, 1,
			version.FirmwareVersion, version.PartNumber, version.SerialNumber)
	}
	ch <- prometheus.MustNewConstMetric(scrapeDuration, prometheus.GaugeValue, verT, apiVersion)

	v, viT, err := scrapeVitals(client)
	if err != nil {
		slog.Error("Failed collection", "err", err, "path", apiVitals)
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0, apiVitals)
	} else {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1, apiVitals)

		var vc float64 = 0
		if v.VehicleConnected {
			vc = 1
		}
		ch <- prometheus.MustNewConstMetric(vehicleConnected, prometheus.GaugeValue, vc)
		ch <- prometheus.MustNewConstMetric(alertsCount, prometheus.GaugeValue, float64(len(v.CurrentAlerts)))
		ch <- prometheus.MustNewConstMetric(gridHz, prometheus.GaugeValue, v.GridHz)
		ch <- prometheus.MustNewConstMetric(gridV, prometheus.GaugeValue, v.GridV)
		ch <- prometheus.MustNewConstMetric(currentAA, prometheus.GaugeValue, v.CurrentAA)
		ch <- prometheus.MustNewConstMetric(currentBA, prometheus.GaugeValue, v.CurrentBA)
		ch <- prometheus.MustNewConstMetric(currentCA, prometheus.GaugeValue, v.CurrentCA)
		ch <- prometheus.MustNewConstMetric(currentNA, prometheus.GaugeValue, v.CurrentNA)
		ch <- prometheus.MustNewConstMetric(voltageAV, prometheus.GaugeValue, v.VoltageAV)
		ch <- prometheus.MustNewConstMetric(voltageBV, prometheus.GaugeValue, v.VoltageBV)
		ch <- prometheus.MustNewConstMetric(voltageCV, prometheus.GaugeValue, v.VoltageCV)
		ch <- prometheus.MustNewConstMetric(handleTemp, prometheus.GaugeValue, v.HandleTempC)
		ch <- prometheus.MustNewConstMetric(mcuTemp, prometheus.GaugeValue, v.McuTempC)
		ch <- prometheus.MustNewConstMetric(pcbaTemp, prometheus.GaugeValue, v.PcbaTempC)
		ch <- prometheus.MustNewConstMetric(sessionEnergyWh, prometheus.CounterValue, v.SessionEnergyWh)
		ch <- prometheus.MustNewConstMetric(sessionSeconds, prometheus.CounterValue, float64(v.SessionS))
		ch <- prometheus.MustNewConstMetric(vehicleCurrentAmps, prometheus.GaugeValue, v.VehicleCurrentA)
	}
	ch <- prometheus.MustNewConstMetric(scrapeDuration, prometheus.GaugeValue, viT, apiVitals)
}

// Tesla Wall Connector Structs
// /api/1/version
type Version struct {
	FirmwareVersion string `json:"firmware_version"`
	PartNumber      string `json:"part_number"`
	SerialNumber    string `json:"serial_number"`
}

func scrapeVersion(client http.Client) (v Version, t float64, err error) {
	start := time.Now()
	resp, err := client.Get(fmt.Sprintf("http://%s%s", *twcAddress, apiVersion))
	if err != nil {
		return v, time.Since(start).Seconds(), err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return v, time.Since(start).Seconds(), err
	}

	t = time.Since(start).Seconds()

	err = json.Unmarshal(body, &v)
	if err != nil {
		return v, t, err
	}

	return v, t, nil
}

// /api/1/lifetime
type Lifetime struct {
	ContactorCycles       int     `json:"contactor_cycles"`
	ContactorCyclesLoaded int     `json:"contactor_cycles_loaded"`
	AlertCount            int     `json:"alert_count"`
	ThermalFoldbacks      int     `json:"thermal_foldbacks"`
	AvgStartupTemp        float64 `json:"avg_startup_temp"`
	ChargeStarts          int     `json:"charge_starts"`
	EnergyWh              int     `json:"energy_wh"`
	ConnectorCycles       int     `json:"connector_cycles"`
	UptimeS               int     `json:"uptime_s"`
	ChargingTimeS         int     `json:"charging_time_s"`
}

func scrapeLifetime(client http.Client) (lt Lifetime, t float64, err error) {
	start := time.Now()
	resp, err := client.Get(fmt.Sprintf("http://%s%s", *twcAddress, apiLifetime))
	if err != nil {
		return lt, time.Since(start).Seconds(), err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return lt, time.Since(start).Seconds(), err
	}
	t = time.Since(start).Seconds()

	err = json.Unmarshal(body, &lt)
	if err != nil {
		return lt, t, err
	}

	return lt, t, nil
}

// /api/1/vitals
type vitals struct {
	ContactorClosed   bool          `json:"contactor_closed"`
	VehicleConnected  bool          `json:"vehicle_connected"`
	SessionS          int           `json:"session_s"`
	GridV             float64       `json:"grid_v"`
	GridHz            float64       `json:"grid_hz"`
	VehicleCurrentA   float64       `json:"vehicle_current_a"`
	CurrentAA         float64       `json:"currentA_a"`
	CurrentBA         float64       `json:"currentB_a"`
	CurrentCA         float64       `json:"currentC_a"`
	CurrentNA         float64       `json:"currentN_a"`
	VoltageAV         float64       `json:"voltageA_v"`
	VoltageBV         float64       `json:"voltageB_v"`
	VoltageCV         float64       `json:"voltageC_v"`
	RelayCoilV        float64       `json:"relay_coil_v"`
	PcbaTempC         float64       `json:"pcba_temp_c"`
	HandleTempC       float64       `json:"handle_temp_c"`
	McuTempC          float64       `json:"mcu_temp_c"`
	UptimeS           int           `json:"uptime_s"`
	InputThermopileUv int           `json:"input_thermopile_uv"`
	ProxV             float64       `json:"prox_v"`
	PilotHighV        float64       `json:"pilot_high_v"`
	PilotLowV         float64       `json:"pilot_low_v"`
	SessionEnergyWh   float64       `json:"session_energy_wh"`
	ConfigStatus      int           `json:"config_status"`
	EvseState         int           `json:"evse_state"`
	CurrentAlerts     []interface{} `json:"current_alerts"`
}

func scrapeVitals(client http.Client) (v vitals, t float64, err error) {
	start := time.Now()
	resp, err := client.Get(fmt.Sprintf("http://%s%s", *twcAddress, apiVitals))
	if err != nil {
		return v, time.Since(start).Seconds(), err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return v, time.Since(start).Seconds(), err
	}
	t = time.Since(start).Seconds()
	err = json.Unmarshal(body, &v)
	if err != nil {
		return v, t, err
	}

	return v, t, nil
}
