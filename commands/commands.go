package commands

import (
	"math/rand"
	"errors"
	"fmt"
	"os"

	"github.com/fotis-sofoulis/pokedex-cli/internal/pokeapi"
	"github.com/fotis-sofoulis/pokedex-cli/internal/pokedex"
)

type Config struct {
	Next     *string
	Previous *string
	LatestEnounters map[string]struct{}
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
		"catch": {
			Name:        "catch <pokemon_name>",
			Description: "Attempt to catch a Pokemon and add it to your Pokedex",
			Callback:    commandCatch,
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
	cfg.LatestEnounters = make(map[string]struct{}) // reset before adding
	for _, encounter := range locationAreaDetails.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
		cfg.LatestEnounters[encounter.Pokemon.Name] = struct{}{}
	}

	return nil
}

func commandCatch(cfg *Config, args ...string) error {
	if len(cfg.LatestEnounters) == 0 {
		return errors.New("You must explore an area before catching Pokémon")
	}

	if len(args) == 0 {
		return errors.New("you must provide a pokemon name")
	}
	name := args[0]

	if _, exist := cfg.LatestEnounters[name]; !exist {
		return fmt.Errorf("%s is not in the currently explored area", name)
	}

	pokemon, rawData, err := pokeapi.GetPokemon(name)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	const catchThreshold = 400
	roll := rand.Intn(catchThreshold)

	if roll < pokemon.BaseExperience {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemon.Name)

	err = pokedex.AddToPokedex(rawData)
	if err != nil {
		return fmt.Errorf("could not add to pokedex: %w", err)
	}

	return nil
}
