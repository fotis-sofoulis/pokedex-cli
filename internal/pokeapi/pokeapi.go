package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fotis-sofoulis/pokedex-cli/internal/pokecache"
)

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaResp struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationAreaDetailsResp struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	ID			   int    `json:"id"`
	Name		   string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

var cache *pokecache.Cache

func InitCache(c *pokecache.Cache) {
	cache = c
}

func FetchLocationAreas(url string) (LocationAreaResp, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	// fetch from cache
	if data, exist := cache.Get(url); exist {
		var res LocationAreaResp
		if err := json.Unmarshal(data, &res); err == nil {
			return res, nil
		}
	}

	res, err := http.Get(url)
	if err != nil {
		return LocationAreaResp{}, fmt.Errorf("failed to fetch location areas: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaResp{}, fmt.Errorf("failed to read response body: %w", err)
	}

	cache.Add(url, body)

	var data LocationAreaResp
	if err := json.Unmarshal(body, &data); err != nil {
		return LocationAreaResp{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data, nil
}

func GetLocationAreaDetails(areaName string) (LocationAreaDetailsResp, error) {
	fullUrl := "https://pokeapi.co/api/v2/location-area/" + areaName

	if data, exist := cache.Get(fullUrl); exist {
		var res LocationAreaDetailsResp
		if err := json.Unmarshal(data, &res); err == nil {
			return res, nil
		}
	}

	res, err := http.Get(fullUrl)
	if err != nil {
		return LocationAreaDetailsResp{}, fmt.Errorf("failed to fetch location area %s: %w", areaName, err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return LocationAreaDetailsResp{}, fmt.Errorf("location area not found")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaDetailsResp{}, fmt.Errorf("failed to read response body: %w", err)
	}

	cache.Add(fullUrl, body)

	var data LocationAreaDetailsResp
	if err := json.Unmarshal(body, &data); err != nil {
		return LocationAreaDetailsResp{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data, nil
}

func GetPokemon(name string) (Pokemon, []byte, error) {
	fullUrl := "https://pokeapi.co/api/v2/pokemon/" + name

	if data, exist := cache.Get(fullUrl); exist {
		var pokemon Pokemon
		if err := json.Unmarshal(data, &pokemon); err == nil {
			return pokemon, data, nil
		}
	}

	res, err := http.Get(fullUrl)
	if err != nil {
		return Pokemon{}, nil, fmt.Errorf("failed to fetch pokemon %s: %w", name, err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return Pokemon{}, nil, fmt.Errorf("pokemon not found")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	cache.Add(fullUrl, body)

	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return Pokemon{}, nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return pokemon, body, nil
}
