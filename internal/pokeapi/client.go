// internal/pokeapi/client.go
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cybergrim/bootdev_pokedex/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	pCache     *pokecache.Cache
}

func NewClient(timeout time.Duration) *Client {
	interval := 5
	return &Client{
		httpClient: http.Client{Timeout: timeout},
		pCache:     pokecache.NewCache(time.Duration(interval) * time.Minute),
	}
}

func (c *Client) ListLocationAreas(pageURL *string) (LocationArea, error) {
	initialURL := "https://pokeapi.co/api/v2/location-area/"
	if pageURL != nil {
		initialURL = *pageURL
	}

	result, exist := c.pCache.Get(initialURL)
	if exist {
		var out LocationArea
		if err := json.Unmarshal(result, &out); err != nil {
			return LocationArea{}, err
		}
		return out, nil
	}

	req, err := http.NewRequest("GET", initialURL, nil)
	if err != nil {
		return LocationArea{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationArea{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return LocationArea{}, fmt.Errorf("pokeapi error: %s: %s", resp.Status, string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationArea{}, err
	}

	c.pCache.Add(initialURL, b)

	var out LocationArea
	if err := json.Unmarshal(b, &out); err != nil {
		return LocationArea{}, err
	}
	return out, nil
}

func (c *Client) ExploreLocation(area string) (ExplorationArea, error) {
	baseURL := "https://pokeapi.co/api/v2/location-area/"
	finalURL := baseURL + area
	result, exist := c.pCache.Get(finalURL)
	if exist {
		var out ExplorationArea
		if err := json.Unmarshal(result, &out); err != nil {
			return ExplorationArea{}, err
		}
		return out, nil
	}

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return ExplorationArea{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ExplorationArea{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return ExplorationArea{}, fmt.Errorf("pokeapi error: %s: %s", resp.Status, string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ExplorationArea{}, err
	}

	c.pCache.Add(finalURL, b)

	var out ExplorationArea
	if err := json.Unmarshal(b, &out); err != nil {
		return ExplorationArea{}, err
	}
	return out, nil
}

func (c *Client) GetPokemonInfo(pokemon string) (PokemonInfo, error) {
	baseURL := "https://pokeapi.co/api/v2/pokemon/"
	finalURL := baseURL + pokemon
	result, exist := c.pCache.Get(finalURL)
	if exist {
		var out PokemonInfo
		if err := json.Unmarshal(result, &out); err != nil {
			return PokemonInfo{}, err
		}
		return out, nil
	}

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return PokemonInfo{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return PokemonInfo{}, fmt.Errorf("pokeapi error: %s: %s", resp.Status, string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonInfo{}, err
	}

	c.pCache.Add(finalURL, b)

	var out PokemonInfo
	if err := json.Unmarshal(b, &out); err != nil {
		return PokemonInfo{}, err
	}
	return out, nil
}
