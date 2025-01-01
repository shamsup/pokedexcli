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

var cache = pokecache.NewCache(5 * time.Minute)

func GetLocations(overrideUrl string) (PaginatedResponse[Location], error) {
	url := baseURL + "location-area/"
	if overrideUrl != "" {
		url = overrideUrl
	}
	result := PaginatedResponse[Location]{}
	if cached, ok := cache.Get(url); ok {
		err := json.Unmarshal(cached, &result)
		if err != nil {
			return result, nil
		}
		// in case of error, we'll just fetch the data again
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return PaginatedResponse[Location]{}, fmt.Errorf("error: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		fmt.Println("Error:", res.Status, string(body))
		return PaginatedResponse[Location]{}, fmt.Errorf("error: %v %v\n%s", res.StatusCode, res.Status, body)
	}
	if err != nil {
		fmt.Println("Error:", err)
		return PaginatedResponse[Location]{}, fmt.Errorf("error: %v", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error:", err)
		return PaginatedResponse[Location]{}, fmt.Errorf("error: %v", err)
	}
	cache.Add(url, body)
	return result, nil
}
