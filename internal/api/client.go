package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"atlas.radio/internal/model"
)

func SearchStations(query string) ([]model.Station, error) {
	// Use the main round-robin address which is more reliable
	baseUrl := "https://all.api.radio-browser.info/json/stations/"
	var apiURL string
	
	if query == "" {
		apiURL = baseUrl + "topclick/50"
	} else {
		apiURL = baseUrl + "byname/" + url.PathEscape(query)
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		// Fallback to a specific stable mirror if round-robin fails
		baseUrl = "https://de1.api.radio-browser.info/json/stations/"
		if query == "" {
			apiURL = baseUrl + "topclick/50"
		} else {
			apiURL = baseUrl + "byname/" + url.PathEscape(query)
		}
		resp, err = http.Get(apiURL)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	var stations []model.Station
	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		return nil, err
	}

	return stations, nil
}
