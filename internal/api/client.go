package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"atlas.radio/internal/model"
)

func SearchStations(query string) ([]model.Station, error) {
	// Use the load-balanced DNS
	baseUrl := "https://at1.api.radio-browser.info/json/stations/"
	var apiURL string
	
	if query == "" {
		apiURL = baseUrl + "topclick/50"
	} else {
		apiURL = baseUrl + "byname/" + url.PathEscape(query)
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stations []model.Station
	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		return nil, err
	}

	return stations, nil
}
