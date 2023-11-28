package pkg

import (
	"log"
	"net/http"
	"crypto/tls"
	"github.com/gcochard/go-envoy"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	descActivePanelCount = prometheus.NewDesc("envoy_active_panel_count", "Number of panels producing power", nil, nil)

	descProductionInverterWatts = prometheus.NewDesc("envoy_production_inverter_watts", "Amount of watts being produced per inverter", nil, nil)
	descProductionInverterWattHours = prometheus.NewDesc("envoy_production_inverter_watthours", "Amount of watthours produced per inverter", nil, nil)
	descProductionWatts = prometheus.NewDesc("envoy_production_watts", "Amount of watts being produced", nil, nil)
	descProductionWattHours = prometheus.NewDesc("envoy_production_watthours", "Amount of watthours produced", nil, nil)
	descProductionWattHoursToday = prometheus.NewDesc("envoy_production_watthours_today", "Amount of watthours produced today", nil, nil)
	descProductionRmsCurrent = prometheus.NewDesc("envoy_production_rms_current_amps", "", nil, nil)
	descProductionRmsVoltage = prometheus.NewDesc("envoy_production_rms_voltage_volts", "", nil, nil)
	descProductionReactivePowerWatts = prometheus.NewDesc("envoy_production_reactive_power_watts", "", nil, nil)
	descProductionApparentPowerWatts = prometheus.NewDesc("envoy_production_apparent_power_watts", "", nil, nil)
	descProductionPowerFactor = prometheus.NewDesc("envoy_production_power_factor", "", nil, nil)

	descConsumptionWatts = prometheus.NewDesc("envoy_consumption_watts", "Amount of watts being consumed", nil, nil)
	descConsumptionWattHours = prometheus.NewDesc("envoy_consumption_watthours", "Amount of watthours consumed", nil, nil)
	descConsumptionWattHoursToday = prometheus.NewDesc("envoy_consumption_watthours_today", "Amount of watthours consumed today", nil, nil)
	descConsumptionRmsCurrent = prometheus.NewDesc("envoy_consumption_rms_current_amps", "", nil, nil)
	descConsumptionRmsVoltage = prometheus.NewDesc("envoy_consumption_rms_voltage_volts", "", nil, nil)
	descConsumptionReactivePowerWatts = prometheus.NewDesc("envoy_consumption_reactive_power_watts", "", nil, nil)
	descConsumptionApparentPowerWatts = prometheus.NewDesc("envoy_consumption_apparent_power_watts", "", nil, nil)
	descConsumptionPowerFactor = prometheus.NewDesc("envoy_consumption_power_factor", "", nil, nil)

	descConsumptionGridWatts = prometheus.NewDesc("envoy_consumption_grid_watts", "Amount of watts being consumed from the grid", nil, nil)
	descConsumptionGridWattHours = prometheus.NewDesc("envoy_consumption_grid_watthours", "Amount of watthours consumed from the grid", nil, nil)
	descConsumptionGridWattHoursToday = prometheus.NewDesc("envoy_consumption_grid_watthours_today", "Amount of watthours consumed from the grid today", nil, nil)
	descConsumptionGridRmsCurrent = prometheus.NewDesc("envoy_consumption_grid_rms_current_amps", "", nil, nil)
	descConsumptionGridRmsVoltage = prometheus.NewDesc("envoy_consumption_grid_rms_voltage_volts", "", nil, nil)
	descConsumptionGridReactivePowerWatts = prometheus.NewDesc("envoy_consumption_grid_reactive_power_watts", "", nil, nil)
	descConsumptionGridApparentPowerWatts = prometheus.NewDesc("envoy_consumption_grid_apparent_power_watts", "", nil, nil)
	descConsumptionGridPowerFactor = prometheus.NewDesc("envoy_consumption_grid_power_factor", "", nil, nil)

	client *envoy.Client
)

type EnvoyCollector struct {
	target string
	token string
}

func NewEnvoyCollector(target string, token string) *EnvoyCollector {
	return &EnvoyCollector{
		target: target,
		token: token,
	}
}

func (s *EnvoyCollector) Describe(chan<- *prometheus.Desc) {
	return
}

func (s *EnvoyCollector) Collect(metrics chan<- prometheus.Metric) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if client == nil {
		log.Printf("Client init")
		client = envoy.NewClientWithHTTP(s.target, "https", &http.Client{Transport: tr})
		if s.token != "" {
			client.SetToken(s.token)
		}
	}

	production, err := client.Production()
	if err != nil {
		log.Printf("failed to get production data from device %s: %s\n", s.target, err.Error())
		return
	}

	for _, section := range production.Production {
		switch(section.Type) {
		case "inverters":
			metrics <- prometheus.MustNewConstMetric(
				descActivePanelCount,
				prometheus.GaugeValue,
				float64(section.ActiveCount),
			)
			metrics <- prometheus.MustNewConstMetric(
				descProductionInverterWattHours,
				prometheus.CounterValue,
				float64(section.WhLifetime),
			)
			metrics <- prometheus.MustNewConstMetric(
				descProductionInverterWatts,
				prometheus.GaugeValue,
				float64(section.WNow),
			)
		case "eim":
			metrics <- prometheus.MustNewConstMetric(
				descProductionRmsCurrent,
				prometheus.GaugeValue,
				float64(section.RmsCurrent),
			)
			metrics <- prometheus.MustNewConstMetric(
				descProductionRmsVoltage,
				prometheus.GaugeValue,
				float64(section.RmsVoltage),
			)
			metrics <- prometheus.MustNewConstMetric(
				descProductionReactivePowerWatts,
				prometheus.GaugeValue,
				float64(section.ReactPwr),
			)
			metrics <- prometheus.MustNewConstMetric(
				descProductionApparentPowerWatts,
				prometheus.GaugeValue,
				float64(section.ApprntPwr),
			)
			metrics <- prometheus.MustNewConstMetric(
				descProductionPowerFactor,
				prometheus.GaugeValue,
				float64(section.PwrFactor),
			)
			metrics <- prometheus.MustNewConstMetric(
				descProductionWattHours,
				prometheus.CounterValue,
				float64(section.WhLifetime),
			)
			metrics <- prometheus.MustNewConstMetric(
				descProductionWattHoursToday,
				prometheus.CounterValue,
				float64(section.WhToday),
			)
			metrics <- prometheus.MustNewConstMetric(
				descProductionWatts,
				prometheus.GaugeValue,
				float64(section.WNow),
			)
		}
	}

	for _, section := range production.Consumption {
		switch (section.MeasurementType) {
		case "total-consumption":
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionWattHours,
				prometheus.CounterValue,
				float64(section.WhLifetime),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionWattHoursToday,
				prometheus.CounterValue,
				float64(section.WhToday),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionWatts,
				prometheus.GaugeValue,
				float64(section.WNow),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionRmsCurrent,
				prometheus.GaugeValue,
				float64(section.RmsCurrent),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionRmsVoltage,
				prometheus.GaugeValue,
				float64(section.RmsVoltage),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionReactivePowerWatts,
				prometheus.GaugeValue,
				float64(section.ReactPwr),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionApparentPowerWatts,
				prometheus.GaugeValue,
				float64(section.ApprntPwr),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionPowerFactor,
				prometheus.GaugeValue,
				float64(section.PwrFactor),
			)
		case "net-consumption":
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionGridWattHours,
				prometheus.CounterValue,
				float64(section.WhLifetime),
			)
			// this metric never seems to change, always 0
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionGridWattHoursToday,
				prometheus.CounterValue,
				float64(section.WhToday),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionGridWatts,
				prometheus.GaugeValue,
				float64(section.WNow),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionGridRmsCurrent,
				prometheus.GaugeValue,
				float64(section.RmsCurrent),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionGridRmsVoltage,
				prometheus.GaugeValue,
				float64(section.RmsVoltage),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionGridReactivePowerWatts,
				prometheus.GaugeValue,
				float64(section.ReactPwr),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionGridApparentPowerWatts,
				prometheus.GaugeValue,
				float64(section.ApprntPwr),
			)
			metrics <- prometheus.MustNewConstMetric(
				descConsumptionGridPowerFactor,
				prometheus.GaugeValue,
				float64(section.PwrFactor),
			)
		}
	}
}
