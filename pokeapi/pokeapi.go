package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ashufelt/pokecache"
)

type LocationConfig struct {
	PreviousLocationPage *string
	NextLocationPage     *string
}

func GetLocationPage(cfg *LocationConfig, cache *pokecache.Cache) error {
	var err error
	var locationResultsPage LocationResultsPage
	cachedResp, exists := cache.Get(*cfg.NextLocationPage)

	if exists { //take from cache
		fmt.Println("Using cached values")
		err = json.Unmarshal(cachedResp, &locationResultsPage)
		if err != nil {
			return err
		}
	} else { //or get data from new GET HTTP request using PokeAPI
		resp, err := http.Get(*cfg.NextLocationPage)
		if err != nil {
			return fmt.Errorf("error reading location %s", *cfg.NextLocationPage)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode > 299 { //if GET returns error code
			return fmt.Errorf("response failed with status code: %d and body: %s", resp.StatusCode, body)
		}
		if err != nil {
			return err
		}
		cache.Add(*cfg.NextLocationPage, []byte(body))
		err = json.Unmarshal([]byte(body), &locationResultsPage)
		if err != nil {
			return err
		}
	}
	err = displayLocationResultsPage(&locationResultsPage)
	if err != nil {
		return err
	}
	cfg.PreviousLocationPage = locationResultsPage.Previous
	cfg.NextLocationPage = locationResultsPage.Next

	for _, location := range locationResultsPage.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func displayLocationResultsPage(locationResultsPage *LocationResultsPage) error {
	for _, location := range locationResultsPage.Results {
		fmt.Printf("%s\n", location.Name)
	}
	return nil
}
