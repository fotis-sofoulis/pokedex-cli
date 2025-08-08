package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"github.com/fotis-sofoulis/pokedex-cli/commands"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &commands.Config{
		Next:     nil,
		Previous: nil,
	}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleaned := cleanInput(input)

		if len(cleaned) == 0 {
			continue
		}

		cmdName := cleaned[0]

		cmd, exists := commands.GetCommands()[cmdName]
		if exists {
			err := cmd.Callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unkown command")
			continue
		}

	}

}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

