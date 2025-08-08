package commands

import (
	"fmt"
	"os"
	"errors"

	"github.com/fotis-sofoulis/pokedex-cli/internal/pokeapi"
)

type Config struct {
	Next     *string
	Previous *string
}

type cliCommand struct {
	Name        string
	Description string
	Callback    func(cfg *Config, args ...string) error
}

func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Displays the next 20 location areas in the Pokemon world",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays the previous 20 location areas in the Pokemon world",
			Callback:    commandMapb,
		},
		"explore": {
			Name:        "explore <location_area>",
			Description: "Lists all the pokemon in a given location area",
			Callback:    commandExplore,
		},
	}
}

func commandExit(cfg *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, args ...string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for _, command := range GetCommands() {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}
	return nil
}

func commandMap(cfg *Config, args ...string) error {
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

func commandMapb(cfg *Config, args ...string) error {
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

func commandExplore(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return errors.New("You must provide a location area name")
	}
	locationAreaName := args[0]
	fmt.Printf("Exploring %s...\n", locationAreaName)

	locationAreaDetails, err := pokeapi.GetLocationAreaDetails(locationAreaName)
	if err != nil {
		return fmt.Errorf("Couldn't get the location are details %w", err)
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationAreaDetails.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
