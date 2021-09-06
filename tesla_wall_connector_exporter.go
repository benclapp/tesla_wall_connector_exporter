package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	uptime = newMetricDesc("uptime_seconds", "Uptime of wallconnector", nil)
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.Info("Tesla Wall Connector Exporter")

	prometheus.MustRegister(NewExporter())
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Tesla Wall Connector Exporter</title></head>
			<body>
			<h1>Tesla Wall Connector Exporter</h1>
			<p><a href=/metrics>Metrics</a></p>
			</body>
			</html>`))
	})

	log.Fatal(http.ListenAndServe(":9851", nil))
}

func newMetricDesc(name string, help string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName("teslawallconnector", "", name), help, labels, nil)
}

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- uptime
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	lt := GetLifetime()
	ch <- prometheus.MustNewConstMetric(uptime, prometheus.GaugeValue, float64(lt.UptimeS))
}

// Tesla Wall Connector Structs
// /api/1/version
type version struct {
	FirmwareVersion string `json:"firmware_version"`
	PartNumber      string `json:"part_number"`
	SerialNumber    string `json:"serial_number"`
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

// TEMP
func GetLifetime() *Lifetime {
	// Perhaps this should return err too, don't fail the entire scrape
	resp, err := http.Get("http://192.168.1.131/api/1/lifetime")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var lt Lifetime

	err = json.Unmarshal(body, &lt)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(lt.UptimeS)

	return &lt
}
