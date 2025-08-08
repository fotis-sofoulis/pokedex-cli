package commands

import (
	"fmt"
	"os"
)

type cliCommand struct {
	Name		string
	Description string
	Callback	func() error
}

func GetCommands() map[string]cliCommand {
	return map[string]cliCommand {
		"exit": {
			Name: "exit",
			Description: "Exit the Pokedex",
			Callback: commandExit,
		},
		"help": {
			Name: "help",
			Description: "Displays a help message",
			Callback: commandHelp,
		},
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for _, command := range GetCommands() {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}
	return nil
}
