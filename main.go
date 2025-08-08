package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var commands = map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleaned := cleanInput(input)

		if len(cleaned) == 0 {
			continue
		}

		switch cleaned[0] {
		case "exit":
			if err := commands["exit"].callback(); err != nil {
				fmt.Println("Error:", err)
			}
		case "help":
			if err := commands["help"].callback(); err != nil {
				fmt.Println("Error:", err)
			}
		default:
			fmt.Println("Unknown command")
		}
	}
}
