package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

var locationCache = map[string]PaginatedResponse[Location]{}

func GetLocations(overrideUrl string) (PaginatedResponse[Location], error) {
	url := baseURL + "location-area/"
	if overrideUrl != "" {
		url = overrideUrl
	}

	if cached, ok := locationCache[url]; ok {
		return cached, nil
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
	result := PaginatedResponse[Location]{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error:", err)
		return PaginatedResponse[Location]{}, fmt.Errorf("error: %v", err)
	}
	locationCache[url] = result
	return result, nil
}
