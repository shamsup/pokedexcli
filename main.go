package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shamsup/pokedexcli/internal/pokeapi"
	"github.com/shamsup/pokedexcli/internal/pokedex"
)

var commands = map[string]Command{}

func registerCommand(cmd Command) {
	commands[cmd.Name] = cmd
}

func main() {
	sharedConfig := Config{Pokedex: pokedex.NewPokedex(pokedex.PokedexConfig{})}

	registerCommand(Command{
		Name:        "help",
		Description: "Displays a help message",
		Handler:     commandHelp,
		Config:      &sharedConfig,
	})
	registerCommand(Command{
		Name:        "exit",
		Description: "Exit the Pokedex",
		Handler:     commandExit,
		Config:      &sharedConfig,
	})

	registerCommand(Command{
		Name:        "map",
		Description: "List locations from the map. Use 'mapb' to go back or 'map' again to go forward",
		Handler:     commandMap,
		Config:      &sharedConfig,
	})

	registerCommand(Command{
		Name:        "mapb",
		Description: "Show the previous page of locations",
		Handler:     commandMapBack,
		Config:      &sharedConfig,
	})

	registerCommand(Command{
		Name:        "explore",
		Description: "Explore a location to find Pokemon",
		Handler:     commandExplore,
	})

	registerCommand(Command{
		Name:        "catch",
		Description: "Catch a Pokemon",
		Handler:     commandCatchPokemon,
		Config:      &sharedConfig,
	})

	registerCommand(Command{
		Name:        "inspect",
		Description: "Inspect a caught Pokemon",
		Handler:     commandInspectPokemon,
		Config:      &sharedConfig,
	})

	registerCommand(Command{
		Name:        "pokedex",
		Description: "List all caught Pokemon",
		Handler:     commandPokedex,
		Config:      &sharedConfig,
	})

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			words := cleanInput(scanner.Text())
			if len(words) > 0 {
				command := words[0]
				args := words[1:]
				if cmd, ok := commands[command]; ok {
					if err := cmd.Handler(cmd.Config, args); err != nil {
						fmt.Println("Error:", err)
					}
				} else {
					fmt.Print("Unknown command\n")
				}
			}
			fmt.Print("Pokedex > ")
		}
	}
}
func cleanInput(text string) []string {
	words := []string{}
	for _, word := range strings.Fields(text) {
		if word != "" {
			words = append(words, strings.ToLower(word))
		}
	}
	return words
}

type Command struct {
	Name        string
	Description string
	Handler     func(c *Config, args []string) error
	Config      *Config
}

type Config struct {
	Next     *string
	Previous *string
	Pokedex  pokedex.Pokedex
}

func commandExit(c *Config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	return nil
}

func commandMap(c *Config, _ []string) error {
	if c.Next == nil && c.Previous != nil {
		fmt.Println("you're on the last page")
		return nil
	}
	if c.Next == nil {
		c.Next = new(string)
	}
	resp, err := pokeapi.GetLocations(*c.Next)
	if err != nil {
		return err
	}
	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}

	c.Next = resp.Next
	c.Previous = resp.Previous
	return nil
}

func commandMapBack(c *Config, _ []string) error {
	if c.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	resp, err := pokeapi.GetLocations(*c.Previous)
	if err != nil {
		return err
	}
	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}

	c.Next = resp.Next
	c.Previous = resp.Previous

	return nil
}

func commandExplore(_ *Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("expected location name")
	}
	location := args[0]
	fmt.Printf("Exploring %s...\n", location)
	details, err := pokeapi.GetLocationDetails(location)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range details.PokemonEncounters {
		fmt.Printf("  - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatchPokemon(c *Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("expected pokemon name")
	}
	pokemon := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
	_, caught, err := c.Pokedex.CatchPokemon(pokemon)
	if err != nil {
		// fmt.Printf("We had trouble finding a %s to catch. Are you sure they're real?", pokemon)
		return err
	}
	if caught {
		fmt.Printf("%s was caught!\n", pokemon)
	} else {
		fmt.Printf("%s got away...\n", pokemon)
	}
	return nil
}

func commandInspectPokemon(c *Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("expected pokemon name")
	}
	name := args[0]
	pokemon, err := c.Pokedex.InspectPokemon(name)
	if err != nil {
		fmt.Println("you have no caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(c *Config, _ []string) error {
	pokemon := c.Pokedex.ListCaughtPokemon()
	for _, p := range pokemon {
		fmt.Printf(" - %s\n", p)
	}
	return nil
}
