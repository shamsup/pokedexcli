package pokedex

import (
	"fmt"
	"testing"

	"github.com/shamsup/pokedexcli/internal/pokeapi"
)

type mockAPIClient struct{}

func (m *mockAPIClient) GetPokemon(name string) (pokeapi.PokemonDetails, error) {
	switch name {
	case "charmander":
		return pokeapi.PokemonDetails{
			Name:           "charmander",
			BaseExperience: 64,
		}, nil
	case "bulbasaur":
		return pokeapi.PokemonDetails{
			Name:           "bulbasaur",
			BaseExperience: 64,
		}, nil
	default:
		return pokeapi.PokemonDetails{}, fmt.Errorf("pokemon not found")
	}
}

func guessTrue(int) bool {
	return true
}

func guessFalse(int) bool {
	return false
}

func TestCatchPokemon(t *testing.T) {

	// Test the CatchPokemon method
	cases := []struct {
		pokemonName string
		expected    bool
		expectedErr error
		config      PokedexConfig
	}{
		{
			pokemonName: "charmander",
			expected:    true,
			expectedErr: nil,
			config: PokedexConfig{
				api:  &mockAPIClient{},
				roll: guessTrue,
			},
		},
		{
			pokemonName: "bulbasaur",
			expected:    false,
			expectedErr: nil,
			config: PokedexConfig{
				api:  &mockAPIClient{},
				roll: guessFalse,
			},
		},
		{
			pokemonName: "pikachu",
			expected:    false,
			expectedErr: fmt.Errorf("pokemon not found"),
			config: PokedexConfig{
				api:  &mockAPIClient{},
				roll: guessTrue,
			},
		},
	}
	for _, c := range cases {
		p := NewPokedex(c.config)
		pokemon, result, err := p.CatchPokemon(c.pokemonName)
		if err != nil && c.expectedErr == nil {
			t.Errorf("unexpected error %v, got %v", c.expectedErr, err)
		}
		if err == nil && c.expectedErr != nil {
			t.Errorf("expected error %v, got %v", c.expectedErr, result)
		}
		if result != c.expected {
			t.Errorf("expected %v, got %v", c.expected, result)
		}
		if result && pokemon.Name != c.pokemonName {
			t.Errorf("expected %v, got %v", c.pokemonName, pokemon.Name)
		}
	}
}
