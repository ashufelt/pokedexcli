package pokeapi

const InitialLocationPage string = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

type LocationResultsPage struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
