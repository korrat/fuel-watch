package main

import "github.com/prometheus/client_golang/prometheus"

var (
	priceDesc = prometheus.NewDesc(
		prometheus.BuildFQName("fuelwatch", "", "price_euros"),
		"The price for the different types of fuel",
		[]string{"stationID", "fuel_name"},
		nil,
	)
	updateTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName("fuelwatch", "", "update_timestamp"),
		"The UNIX timestamp when the prices were last updated",
		[]string{"stationID"},
		nil,
	)
)

type fuelCollector struct{}

func (fuelCollector) Describe(desc chan<- *prometheus.Desc) {
	desc <- priceDesc
	desc <- updateTimeDesc
}

func (fuelCollector) Collect(ch chan<- prometheus.Metric) {
	fi := readFuelInfo()

	for name, p := range fi.prices {
		ch <- prometheus.MustNewConstMetric(priceDesc, prometheus.GaugeValue, p, fi.station, name)
	}

	ch <- prometheus.MustNewConstMetric(updateTimeDesc, prometheus.GaugeValue, float64(fi.lastUpdate.Unix()),
		fi.station)
}
