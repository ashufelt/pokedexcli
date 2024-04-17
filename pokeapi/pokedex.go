package pokeapi

import (
	"sync"
)

type Pokedex struct {
	mu            *sync.Mutex
	PokemonCaught map[string]PokemonEntry
}

type PokemonEntry struct {
	val []byte
}

func NewPokedex() Pokedex {
	c := Pokedex{PokemonCaught: make(map[string]PokemonEntry), mu: &sync.Mutex{}}
	return c
}

func (p *Pokedex) Add(key string, pokemonBytes []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.PokemonCaught[key] = PokemonEntry{val: pokemonBytes}
}

func (p *Pokedex) Get(key string) ([]byte, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if data, ok := p.PokemonCaught[key]; ok {
		return data.val, true
	} else {
		return nil, false
	}
}
