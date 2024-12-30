package main

import (
	"bufio"
	"fmt"
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

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			words := cleanInput(scanner.Text())
			if len(words) > 0 {
				command := words[0]
				if cmd, ok := commands[command]; ok {
					if err := cmd.Handler(); err != nil {
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
	Handler     func() error
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	return nil
}
