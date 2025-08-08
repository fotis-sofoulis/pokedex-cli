package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`

}

type LocationAreaResp struct {
	Count    int      		 `json:"count"`
	Next     *string  		 `json:"next"`
	Previous *string  		 `json:"previous"`
	Results  []LocationArea  `json:"results"`
}

func FetchLocationAreas(url string) (LocationAreaResp, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	res, err := http.Get(url)
	if err != nil {
		return LocationAreaResp{}, fmt.Errorf("failed to fetch location areas: %w", err)
	}
	defer res.Body.Close()

	var data LocationAreaResp
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return LocationAreaResp{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data, nil
} 
