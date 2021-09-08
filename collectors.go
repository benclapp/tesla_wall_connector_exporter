package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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
	ch <- uptime
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Scrape the Wall Connector, unmarshal to structs. Update up metric with status
	lt, err := scrapeLifetime()
	if err != nil {
		log.Error(err)
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0, apiLifetime)
	} else {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1, apiLifetime)

		ch <- prometheus.MustNewConstMetric(chargeStarts, prometheus.GaugeValue, float64(lt.ChargeStarts))
		ch <- prometheus.MustNewConstMetric(chargingTime, prometheus.GaugeValue, float64(lt.ChargingTimeS))
		ch <- prometheus.MustNewConstMetric(energyWh, prometheus.GaugeValue, float64(lt.EnergyWh))
		ch <- prometheus.MustNewConstMetric(uptime, prometheus.GaugeValue, float64(lt.UptimeS))
	}

	version, err := scrapeVersion()
	if err != nil {
		log.Error(err)
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0, apiVersion)
	} else {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1, apiVersion)

		ch <- prometheus.MustNewConstMetric(info, prometheus.GaugeValue, 1,
			version.FirmwareVersion, version.PartNumber, version.SerialNumber)
	}
	v, err := scrapeVitals()
	if err != nil {
		log.Error(err)
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
		ch <- prometheus.MustNewConstMetric(handleTemp, prometheus.GaugeValue, v.HandleTempC)
		ch <- prometheus.MustNewConstMetric(mcuTemp, prometheus.GaugeValue, v.McuTempC)
		ch <- prometheus.MustNewConstMetric(pcbaTemp, prometheus.GaugeValue, v.PcbaTempC)
		ch <- prometheus.MustNewConstMetric(sessionEnergyWh, prometheus.CounterValue, v.SessionEnergyWh)
		ch <- prometheus.MustNewConstMetric(sessionSeconds, prometheus.CounterValue, float64(v.SessionS))
		ch <- prometheus.MustNewConstMetric(vehicleCurrentAmps, prometheus.GaugeValue, v.VehicleCurrentA)
	}

}

// Tesla Wall Connector Structs
// /api/1/version
type Version struct {
	FirmwareVersion string `json:"firmware_version"`
	PartNumber      string `json:"part_number"`
	SerialNumber    string `json:"serial_number"`
}

func scrapeVersion() (v Version, err error) {
	resp, err := http.Get(fmt.Sprintf("http://%s%s", *twcAddress, apiVersion))
	if err != nil {
		log.Debug(err)
		return v, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug(err)
		return v, err
	}

	log.Debug(body)
	err = json.Unmarshal(body, &v)
	if err != nil {
		log.Debug(err)
		return v, err
	}

	return v, nil
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

func scrapeLifetime() (lt Lifetime, err error) {
	resp, err := http.Get(fmt.Sprintf("http://%s%s", *twcAddress, apiLifetime))
	if err != nil {
		log.Debug(err)
		return lt, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug(err)
		return lt, err
	}

	log.Debug(body)
	err = json.Unmarshal(body, &lt)
	if err != nil {
		log.Debug(err)
		return lt, err
	}

	return lt, nil
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

func scrapeVitals() (v vitals, err error) {
	resp, err := http.Get(fmt.Sprintf("http://%s%s", *twcAddress, apiVitals))
	if err != nil {
		log.Debug(err)
		return v, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug(err)
		return v, err
	}

	log.Debug(body)
	err = json.Unmarshal(body, &v)
	if err != nil {
		log.Debug(err)
		return v, err
	}

	return v, nil
}
