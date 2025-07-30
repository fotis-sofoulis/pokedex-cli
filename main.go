package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func main() {
	fmt.Println("Hello, World!")
}
