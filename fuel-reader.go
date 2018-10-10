package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	stationID = os.Getenv("PFR_STATION_ID")
)

type fuelInfo struct {
	station    string
	prices     map[string]float64
	lastUpdate time.Time
}

func (fi *fuelInfo) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v == nil {
		return nil
	}

	// Runtime panic possible, should be protected with recover
	fieldMap := v.(map[string]interface{})["response"].(map[string]interface{})

	id := fieldMap["stationId"].(string)

	lastUpdate, err := time.Parse("200601021504", fieldMap["lastUpdate"].(string))
	if err != nil {
		return err
	}

	prices := make(map[string]float64)
	for _, p := range fieldMap["prices"].([]interface{}) {
		m := p.(map[string]interface{})
		prices[m["name"].(string)], err = strconv.ParseFloat(strings.Replace(m["price"].(string), ",", ".", -1), 64)
		if err != nil {
			return err
		}
	}

	fi.station = id
	fi.lastUpdate = lastUpdate.In(time.UTC)
	fi.prices = prices

	return nil
}

func readFuelInfo() fuelInfo {
	v := url.Values{}
	v.Add("stationId", stationID)

	resp, err := http.Get("https://ap.aral.de/api/v2/getStationPricesById.php?" + v.Encode())
	if err != nil {
		panic(err)
	}

	d := json.NewDecoder(resp.Body)

	var r fuelInfo
	err = d.Decode(&r)
	if err != nil {
		panic(err)
	}

	return r
}
