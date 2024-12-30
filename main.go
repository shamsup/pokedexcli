package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			words := cleanInput(scanner.Text())
			if len(words) > 0 {
				command := words[0]
				fmt.Printf("Your command was: %s\n", command)
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
