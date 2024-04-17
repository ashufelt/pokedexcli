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

func GetLocationsDump(cfg *LocationConfig, cache *pokecache.Cache) error {
	var err error
	var locationResultsPage LocationDumpResults = LocationDumpResults{}
	cachedResp, exists := cache.Get(*cfg.NextLocationPage)

	if exists { //take from cache
		//fmt.Println("Using cached values")
		err = json.Unmarshal(cachedResp, &locationResultsPage)
		if err != nil {
			return err
		}
	} else { //or get data from new GET HTTP request using PokeAPI
		err = fillDataFromEndpoint(*cfg.NextLocationPage, cache, &locationResultsPage)
		if err != nil {
			return err
		}
	}
	displayLocationDumpResultsPage(&locationResultsPage)
	cfg.PreviousLocationPage = locationResultsPage.Previous
	cfg.NextLocationPage = locationResultsPage.Next

	return nil
}

func displayLocationDumpResultsPage(locationResultsPage *LocationDumpResults) {
	for _, location := range locationResultsPage.Results {
		fmt.Printf("%s\n", location.Name)
	}
}

func GetSpecificLocationInfo(cache *pokecache.Cache, areaName string) error {
	var err error
	endpoint := fmt.Sprintf("%s%s", BaseLocationEndpoint, areaName)
	var locationAreaInfo LocationAreaInformation = LocationAreaInformation{}

	cachedResp, exists := cache.Get(endpoint)
	if exists { //take from cache
		//fmt.Println("Using cached values")
		err = json.Unmarshal(cachedResp, &locationAreaInfo)
		if err != nil {
			return err
		}
	} else { //or get data from new GET HTTP request using PokeAPI
		err = fillDataFromEndpoint(endpoint, cache, &locationAreaInfo)
		if err != nil {
			fmt.Printf("Could not receive data for location '%s'\n", areaName)
			return err
		}
	}
	displayLocationAreaInformation(&locationAreaInfo)
	return nil
}

func displayLocationAreaInformation(locationAreaInformation *LocationAreaInformation) {
	if len(locationAreaInformation.PokemonEncounters) == 0 {
		fmt.Println("No Pokemon found, be sure you typed the location area name correctly")
	}
	fmt.Println("Found Pokemon:")
	for _, pokemon_en := range locationAreaInformation.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon_en.Pokemon.Name)
	}
}

func fillDataFromEndpoint(endpoint string, cache *pokecache.Cache, locationDataStruct LocationStruct) error {
	resp, err := http.Get(endpoint)
	if err != nil {
		return fmt.Errorf("error with GET endpoint %s", endpoint)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 { //if GET returns error code
		return fmt.Errorf("response failed with status code: %d and body: %s", resp.StatusCode, body)
	}
	if err != nil {
		return err
	}
	cache.Add(endpoint, []byte(body))
	err = json.Unmarshal([]byte(body), locationDataStruct)
	if err != nil {
		return err
	}
	return nil
}
