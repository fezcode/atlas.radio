package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"atlas.radio/internal/model"
)

func SearchStations(query string) ([]model.Station, error) {
	// Use a highly stable mirror
	baseUrl := "https://de1.api.radio-browser.info/json/stations/"
	var apiURL string
	
	if query == "" {
		// Filter for MP3 and high bitrate to avoid decoding issues
		apiURL = baseUrl + "topclick/50?codec=MP3"
	} else {
		apiURL = baseUrl + "byname/" + url.PathEscape(query) + "?codec=MP3"
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
