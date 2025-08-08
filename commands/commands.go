package commands

import (
	"fmt"
	"os"
	"github.com/fotis-sofoulis/pokedex-cli/internal/pokeapi"
)

type Config struct {
	Next	 *string
	Previous *string
}

type cliCommand struct {
	Name		string
	Description string
	Callback	func(*Config) error
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
		"map": {
			Name: "map",
			Description: "Displays the next 20 location areas in the Pokemon world",
			Callback: commandMap,
		},
		"mapb": {
			Name: "mapb",
			Description: "Displays the previous 20 location areas in the Pokemon world",
			Callback: commandMapb,
		},

	}
}

func commandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for _, command := range GetCommands() {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}
	return nil
}

func commandMap(cfg *Config) error {
	url := ""
	if cfg.Next != nil {
		url = *cfg.Next
	}

	data, err := pokeapi.FetchLocationAreas(url)
	if err != nil {
		return err
	}

	for _, loc := range data.Results {
		fmt.Println(loc.Name)
	}

	cfg.Next = data.Next
	cfg.Previous = data.Previous

	return nil
}

func commandMapb(cfg * Config) error {
	if cfg.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	data, err := pokeapi.FetchLocationAreas(*cfg.Previous)
    if err != nil {
        return err
    }

    for _, loc := range data.Results {
        fmt.Println(loc.Name)
    }

    cfg.Next = data.Next
    cfg.Previous = data.Previous

    return nil
}
