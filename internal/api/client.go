package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"atlas.radio/internal/model"
)

func SearchStations(query string) ([]model.Station, error) {
	// Using a reliable Radio Browser mirror
	apiURL := "https://de1.api.radio-browser.info/json/stations/byname/" + url.PathEscape(query)
	if query == "" {
		// Default to popular/top-clicked if no query
		apiURL = "https://de1.api.radio-browser.info/json/stations/topclick/20"
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

func SearchByTag(tag string) ([]model.Station, error) {
	apiURL := "https://de1.api.radio-browser.info/json/stations/bytag/" + url.PathEscape(tag)
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
