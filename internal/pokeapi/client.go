// internal/pokeapi/client.go
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{httpClient: http.Client{Timeout: timeout}}
}

func (c *Client) ListLocationAreas(pageURL *string) (LocationArea, error) {
	initialURL := "https://pokeapi.co/api/v2/location-area/"
	if pageURL != nil {
		initialURL = *pageURL
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

	var out LocationArea
	if err := json.Unmarshal(b, &out); err != nil {
		return LocationArea{}, err
	}
	return out, nil
}
