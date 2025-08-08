package main

import (
	"time"

	"github.com/fotis-sofoulis/pokedex-cli/internal/pokeapi"
	"github.com/fotis-sofoulis/pokedex-cli/internal/pokecache"
)

func main() {
	cache := pokecache.NewCache(5 * time.Second)
	pokeapi.InitCache(cache)
	startRepl()
	}
