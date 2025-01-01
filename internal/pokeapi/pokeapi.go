package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shamsup/pokedexcli/internal/pokecache"
)

var baseURL = "https://pokeapi.co/api/v2/"

type PaginatedResponse[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationDetails struct {
	ID                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	GameIndex            int                    `json:"game_index"`
	EncounterMethodRates []EncounterMethodRates `json:"encounter_method_rates"`
	Location             Location               `json:"location"`
	Names                []LocationNames        `json:"names"`
	PokemonEncounters    []PokemonEncounters    `json:"pokemon_encounters"`
}
type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type VersionDetails struct {
	Rate    int     `json:"rate"`
	Version Version `json:"version"`
}
type EncounterMethodRates struct {
	EncounterMethod EncounterMethod  `json:"encounter_method"`
	VersionDetails  []VersionDetails `json:"version_details"`
}
type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type LocationNames struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}
type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Method struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type EncounterDetails struct {
	MinLevel        int    `json:"min_level"`
	MaxLevel        int    `json:"max_level"`
	ConditionValues []any  `json:"condition_values"`
	Chance          int    `json:"chance"`
	Method          Method `json:"method"`
}
type EncounterVersionDetails struct {
	Version          Version            `json:"version"`
	MaxChance        int                `json:"max_chance"`
	EncounterDetails []EncounterDetails `json:"encounter_details"`
}
type PokemonEncounters struct {
	Pokemon        Pokemon                   `json:"pokemon"`
	VersionDetails []EncounterVersionDetails `json:"version_details"`
}

var cache = pokecache.NewCache(5 * time.Minute)

func GetLocations(overrideUrl string) (PaginatedResponse[Location], error) {
	url := baseURL + "location-area/"
	if overrideUrl != "" {
		url = overrideUrl
	}
	result, err := cachedFetch[PaginatedResponse[Location]](url)
	return result, err
}

func GetLocationDetails(location string) (LocationDetails, error) {
	url := baseURL + "location-area/" + location
	result, err := cachedFetch[LocationDetails](url)
	return result, err
}

func cachedFetch[Response any](url string) (Response, error) {
	var result Response
	var zero Response
	if cached, ok := cache.Get(url); ok {
		err := json.Unmarshal(cached, &result)
		if err == nil {
			return result, nil
		}
		// in case of error, we'll just fetch the data again
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return zero, fmt.Errorf("error: %v", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		fmt.Println("Error:", res.Status, string(resBody))
		return zero, fmt.Errorf("error: %v %v\n%s", res.StatusCode, res.Status, resBody)
	}

	if err != nil {
		fmt.Println("Error:", err)
		return zero, fmt.Errorf("error: %v", err)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		fmt.Println("Error:", err)
		return zero, fmt.Errorf("error: %v", err)
	}
	cache.Add(url, resBody)
	return result, nil
}
