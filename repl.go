package main

import (
	"bufio"
	"fmt"
	"github.com/fotis-sofoulis/pokedex-cli/commands"
	"os"
	"strings"
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
		args := []string{}
		if len(cleaned) > 1 {
			args = cleaned[1:]
		}

		cmd, exists := commands.GetCommands()[cmdName]
		if exists {
			err := cmd.Callback(cfg, args...)
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
