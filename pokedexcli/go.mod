module github.com/ashufelt/pokedexcli

go 1.22.2

replace github.com/ashufelt/pokeapi v0.0.0 => ../pokeapi

replace github.com/ashufelt/pokecache v0.0.0 => ../pokecache

require (
	github.com/ashufelt/pokeapi v0.0.0
	github.com/ashufelt/pokecache v0.0.0
)
