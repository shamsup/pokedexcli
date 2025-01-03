package pokedex

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/shamsup/pokedexcli/internal/pokeapi"
)

type pokedexEntry struct {
	Pokemon   pokeapi.PokemonDetails
	Collected bool
}

type Pokedex struct {
	collection map[string]pokedexEntry
	api        APIClient
	roll       func(int) bool
}

type APIClient interface {
	GetPokemon(name string) (pokeapi.PokemonDetails, error)
}

func (p *Pokedex) SeenPokemon(name string) bool {
	if pokemon, ok := p.collection[name]; ok {
		return pokemon.Collected
	}
	return false
}

func roll(baseExperience int) bool {
	// fmt.Printf("Base experience: %d\n", baseExperience)
	odds := int64(math.Max(math.Sqrt(math.Max(1.0, float64(baseExperience-40))), 1.0))
	// fmt.Printf("Odds of catching: 1/%d\n", odds)
	return (rand.Int63n(odds) + 1) == 1
}

func (p *Pokedex) CatchPokemon(name string) (pokeapi.PokemonDetails, bool, error) {
	var zeroPokemon pokeapi.PokemonDetails
	if pokemon, ok := p.collection[name]; ok {
		collected := p.roll(pokemon.Pokemon.BaseExperience)

		pokemon.Collected = collected
		p.collection[name] = pokemon
		return pokemon.Pokemon, collected, nil
	}
	pokemon, err := p.api.GetPokemon(name)
	if err != nil {
		return zeroPokemon, false, err
	}
	collected := p.roll(pokemon.BaseExperience)
	p.collection[name] = pokedexEntry{
		Pokemon:   pokemon,
		Collected: collected,
	}
	return pokemon, collected, nil
}

func (p *Pokedex) InspectPokemon(name string) (pokeapi.PokemonDetails, error) {
	var zeroPokemon pokeapi.PokemonDetails
	if pokemon, ok := p.collection[name]; ok && pokemon.Collected {
		return pokemon.Pokemon, nil
	}
	return zeroPokemon, fmt.Errorf("you have not caught that pokemon")
}

func (p *Pokedex) ListCaughtPokemon() []string {
	var collected []string
	for name, entry := range p.collection {
		if entry.Collected {
			collected = append(collected, name)
		}
	}
	return collected
}

type DefaultAPIClient struct{}

func (DefaultAPIClient) GetPokemon(name string) (pokeapi.PokemonDetails, error) {
	return pokeapi.GetPokemon(name)
}

type PokedexConfig struct {
	api  APIClient
	roll func(int) bool
}

func NewPokedex(config PokedexConfig) Pokedex {
	if config.api == nil {
		config.api = DefaultAPIClient{}
	}
	if config.roll == nil {
		config.roll = roll
	}
	collection := make(map[string]pokedexEntry)
	return Pokedex{collection: collection, api: config.api, roll: config.roll}
}
