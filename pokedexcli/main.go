package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ashufelt/pokeapi"
	"github.com/ashufelt/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.LocationConfig, *pokecache.Cache, *pokeapi.Pokedex, string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"catch": {
			name:        "catch",
			description: "with <pokemon-name>, attempt to catch that pokemon",
			callback:    commandCatch,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore",
			description: "with <area-name>, get information about an area",
			callback:    commandExplore,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays info about a Pokemon, if you have caught it",
			callback:    commandInspect,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations",
			callback:    commandMapB,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Display all pokemon you have caught",
			callback:    commandPokedex,
		},
	}
}

func commandCatch(cfg *pokeapi.LocationConfig, cache *pokecache.Cache, myDex *pokeapi.Pokedex, pokemonName string) error {
	fmt.Printf("Throwing a Pokeball at %s\n", pokemonName)
	return pokeapi.CatchPokemonAttempt(cache, myDex, pokemonName)
}

func commandExit(cfg *pokeapi.LocationConfig, cache *pokecache.Cache, _ *pokeapi.Pokedex, _ string) error {
	os.Exit(0)
	return nil
}

func commandExplore(cfg *pokeapi.LocationConfig, cache *pokecache.Cache, _ *pokeapi.Pokedex, areaName string) error {
	fmt.Printf("--- Exploring %s ---\n", areaName)
	return pokeapi.GetSpecificLocationInfo(cache, areaName)
}

func commandHelp(cfg *pokeapi.LocationConfig, cache *pokecache.Cache, _ *pokeapi.Pokedex, _ string) error {
	commands := getCommands()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("---------Usage---------")
	fmt.Println("")
	for _, c := range commands {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	fmt.Println("")
	return nil
}

func commandInspect(_ *pokeapi.LocationConfig, _ *pokecache.Cache, myDex *pokeapi.Pokedex, pokemonName string) error {
	var pokemon pokeapi.Pokemon = pokeapi.Pokemon{}

	dexInfo, exists := myDex.Get(pokemonName)
	if !exists {
		fmt.Println("You have not caught this pokemon")
		return nil
	}
	json.Unmarshal(dexInfo, &pokemon)
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%s: %v\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Type(s):")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}
	return nil
}

func commandMap(cfg *pokeapi.LocationConfig, cache *pokecache.Cache, _ *pokeapi.Pokedex, _ string) error {
	fmt.Println("------Retrieving up to 20 location areas------")
	return pokeapi.GetLocationsDump(cfg, cache)
}

func commandMapB(cfg *pokeapi.LocationConfig, cache *pokecache.Cache, _ *pokeapi.Pokedex, _ string) error {
	if cfg.PreviousLocationPage == nil {
		fmt.Println("Already on the first page")
		return nil
	}
	cfg.NextLocationPage = copystring(*cfg.PreviousLocationPage)
	return commandMap(cfg, cache, nil, "")
}

func commandPokedex(_ *pokeapi.LocationConfig, _ *pokecache.Cache, myDex *pokeapi.Pokedex, _ string) error {
	fmt.Println("Your Pokedex:")
	for key := range myDex.PokemonCaught {
		fmt.Printf(" - %s\n", key)
	}
	return nil
}

func copystring(a string) *string {
	if len(a) == 0 {
		return nil
	}
	b := a[0:1] + a[1:]
	return &b
}

func main() {
	configuration := pokeapi.LocationConfig{
		NextLocationPage:     copystring(pokeapi.InitialLocationPage),
		PreviousLocationPage: nil,
	}
	pokeCache := pokecache.NewCache(30 * time.Second)
	myPokedex := pokeapi.NewPokedex()
	inputScanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	for {
		fmt.Print("pokedexcli > ")
		inputScanner.Scan()
		input := inputScanner.Text()
		splitInput := strings.Split(input, " ")
		if len(input) == 0 {
			continue
		}
		if c, ok := commands[splitInput[0]]; ok {
			if len(splitInput) == 1 {
				c.callback(&configuration, &pokeCache, &myPokedex, "")
			} else if len(splitInput) == 2 {
				c.callback(&configuration, &pokeCache, &myPokedex, splitInput[1])
			}

		}
	}
}
