//internal/pokeapi/types.go

package pokeapi

type PokeAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationArea struct {
	Count   int               `json:"count"`
	Next    *string           `json:"next"`
	Prev    *string           `json:"previous"`
	Results []PokeAPIResource `json:"results"`
}
