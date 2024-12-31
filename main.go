package main

import (
	"bufio"
	"fmt"
	"internal/pokeapi"
	"os"
	"strings"
)

var commands = map[string]Command{}

func registerCommand(cmd Command) {
	commands[cmd.Name] = cmd
}

func main() {
	registerCommand(Command{
		Name:        "help",
		Description: "Displays a help message",
		Handler:     commandHelp,
	})
	registerCommand(Command{
		Name:        "exit",
		Description: "Exit the Pokedex",
		Handler:     commandExit,
	})

	mapConfig := Config{}
	registerCommand(Command{
		Name:        "map",
		Description: "List locations from the map. Use 'mapb' to go back or 'map' again to go forward",
		Handler:     commandMap,
		Config:      &mapConfig,
	})

	registerCommand(Command{
		Name:        "mapb",
		Description: "Show the previous page of locations",
		Handler:     commandMapBack,
		Config:      &mapConfig,
	})

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			words := cleanInput(scanner.Text())
			if len(words) > 0 {
				command := words[0]
				if cmd, ok := commands[command]; ok {
					if err := cmd.Handler(cmd.Config); err != nil {
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
	Handler     func(c *Config) error
	Config      *Config
}

type Config struct {
	Next     *string
	Previous *string
}

func commandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	return nil
}

func commandMap(c *Config) error {
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

func commandMapBack(c *Config) error {
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
