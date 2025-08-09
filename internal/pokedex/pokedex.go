package pokedex

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.oneofone.dev/resize"
)

const (
	alphaThreshold = 32768
	reset          = "\x1b[0m"
	cacheDir       = ".cache"
	spriteWidth    = 48
	spriteHeight   = 48
)

var TypeColorMap = map[string]string{
	"normal":   "\x1b[48;2;168;168;120m\x1b[38;2;255;255;255m Normal \x1b[0m",
	"fire":     "\x1b[48;2;240;128;48m\x1b[38;2;255;255;255m Fire \x1b[0m",
	"water":    "\x1b[48;2;104;144;240m\x1b[38;2;255;255;255m Water \x1b[0m",
	"electric": "\x1b[48;2;248;208;48m\x1b[38;2;255;255;255m Electric \x1b[0m",
	"grass":    "\x1b[48;2;120;200;80m\x1b[38;2;255;255;255m Grass \x1b[0m",
	"ice":      "\x1b[48;2;152;216;216m\x1b[38;2;255;255;255m Ice \x1b[0m",
	"fighting": "\x1b[48;2;192;48;40m\x1b[38;2;255;255;255m Fighting \x1b[0m",
	"poison":   "\x1b[48;2;160;64;160m\x1b[38;2;255;255;255m Poison \x1b[0m",
	"ground":   "\x1b[48;2;224;192;104m\x1b[38;2;255;255;255m Ground \x1b[0m",
	"flying":   "\x1b[48;2;168;144;240m\x1b[38;2;255;255;255m Flying \x1b[0m",
	"psychic":  "\x1b[48;2;248;88;136m\x1b[38;2;255;255;255m Psychic \x1b[0m",
	"bug":      "\x1b[48;2;168;184;32m\x1b[38;2;255;255;255m Bug \x1b[0m",
	"rock":     "\x1b[48;2;184;160;56m\x1b[38;2;255;255;255m Rock \x1b[0m",
	"ghost":    "\x1b[48;2;112;88;152m\x1b[38;2;255;255;255m Ghost \x1b[0m",
	"dragon":   "\x1b[48;2;112;56;248m\x1b[38;2;255;255;255m Dragon \x1b[0m",
	"dark":     "\x1b[48;2;112;88;72m\x1b[38;2;255;255;255m Dark \x1b[0m",
	"steel":    "\x1b[48;2;184;184;208m\x1b[38;2;255;255;255m Steel \x1b[0m",
	"fairy":    "\x1b[48;2;238;153;172m\x1b[38;2;255;255;255m Fairy \x1b[0m",
}

func AddToPokedex(pokemonDataRaw []byte) error {
	var pokeData map[string]any
	if err := json.Unmarshal(pokemonDataRaw, &pokeData); err != nil {
		return fmt.Errorf("failed to parse pokemon data: %w", err)
	}

	name := pokeData["name"].(string)
	if ok, _ := IsCaught(name); ok {
		fmt.Printf("%s is already in your Pokedex!\n", name)
		return nil
	}

	pokeID, err := renderPokemonFromData(pokeData)
	if err != nil {
		return fmt.Errorf("failed to process %s: %w", name, err)
	}

	if err := addCaught(pokeID, name); err != nil {
		return fmt.Errorf("failed to save %s to pokedex: %w", name, err)
	}
	return nil
}

func formatTypes(types []string) string {
	var result strings.Builder
	for i, t := range types {
		result.WriteString(TypeColorMap[t])
		if i < len(types)-1 {
			result.WriteString(" | ")
		}
	}
	return result.String()
}

func renderPokemonFromData(pokeData map[string]any) (int, error) {
	pokeID := int(pokeData["id"].(float64))

	// Get stats
	stats := pokeData["stats"].([]any)
	hp := int(stats[0].(map[string]any)["base_stat"].(float64))
	attack := int(stats[1].(map[string]any)["base_stat"].(float64))
	defense := int(stats[2].(map[string]any)["base_stat"].(float64))
	spAtk := int(stats[3].(map[string]any)["base_stat"].(float64))
	spDef := int(stats[4].(map[string]any)["base_stat"].(float64))
	speed := int(stats[5].(map[string]any)["base_stat"].(float64))

	typesData, ok := pokeData["types"].([]any)
	if !ok {
		return 0, fmt.Errorf("'types' field is missing or not an array")
	}

	types := make([]string, 0, 2)
	types = append(types, typesData[0].(map[string]any)["type"].(map[string]any)["name"].(string))
	if len(typesData) > 1 {
		types = append(types, typesData[1].(map[string]any)["type"].(map[string]any)["name"].(string))
	}

	height := pokeData["height"].(float64) / 10.0
	weight := pokeData["weight"].(float64) / 10.0

	pokeInfo := []string{
		"\x1b[47m\x1b[30m═════════ POKÉDEX DATA ═════════\x1b[0m",
		fmt.Sprintf("\x1b[1mName:\x1b[0m     %s", pokeData["name"].(string)),
		fmt.Sprintf("\x1b[1mID:\x1b[0m       #%d", pokeID),
		fmt.Sprintf("\x1b[1mType:\x1b[0m     %s", formatTypes(types)),
		fmt.Sprintf("\x1b[1mHeight:\x1b[0m   %.2f m", height),
		fmt.Sprintf("\x1b[1mWeight:\x1b[0m   %.1f kg", weight),
		"",
		"\x1b[47m\x1b[30m═════════ BASE STATS ══════════\x1b[0m",
		fmt.Sprintf("\x1b[1mHP:\x1b[0m       %d", hp),
		fmt.Sprintf("\x1b[1mAttack:\x1b[0m   %d", attack),
		fmt.Sprintf("\x1b[1mDefense:\x1b[0m  %d", defense),
		fmt.Sprintf("\x1b[1mSp.Atk:\x1b[0m   %d", spAtk),
		fmt.Sprintf("\x1b[1mSp.Def:\x1b[0m   %d", spDef),
		fmt.Sprintf("\x1b[1mSpeed:\x1b[0m    %d", speed),
	}

	spriteURL, _ := pokeData["sprites"].(map[string]any)["front_default"].(string)
	ascii_sprite, err := imageToAscii(spriteURL)
	if err != nil {
		return 0, err
	}

	var combined []string
	maxLines := max(len(ascii_sprite), len(pokeInfo))
	for i := 0; i < maxLines; i++ {
		var left, right string
		if i < len(ascii_sprite) {
			left = ascii_sprite[i]
		}
		if i < len(pokeInfo) {
			right = pokeInfo[i]
		}
		combined = append(combined, fmt.Sprintf("%-*s  %s", 52, left, right))
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return 0, err
	}

	cachePath := filepath.Join(cacheDir, fmt.Sprintf("%s.txt", pokeData["name"].(string)))
	if err := os.WriteFile(cachePath, []byte(strings.Join(combined, "\n")+"\n"), 0644); err != nil {
		return 0, err
	}

	return pokeID, nil
}

func imageToAscii(url string) ([]string, error) {
	sprite, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sprite: %w", err)
	}
	defer sprite.Body.Close()

	img, _, err := image.Decode(sprite.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	resized := resize.Resize(spriteWidth, spriteHeight, img, resize.NearestNeighbor)
	bounds := resized.Bounds()

	var spriteLines []string
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		var line string
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			topR, topG, topB, topA := resized.At(x, y).RGBA()

			var botR, botG, botB, botA uint32
			if y+1 < bounds.Max.Y {
				botR, botG, botB, botA = resized.At(x, y+1).RGBA()
			}

			switch {
			case topA > alphaThreshold && botA > alphaThreshold:
				line += fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", topR>>8, topG>>8, topB>>8, botR>>8, botG>>8, botB>>8)
			case topA > alphaThreshold:
				line += fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[49m▀", topR>>8, topG>>8, topB>>8)
			case botA > alphaThreshold:
				line += fmt.Sprintf("\x1b[48;2;%d;%d;%dm▄", botR>>8, botG>>8, botB>>8)
			default:
				line += reset + " "
			}
		}
		line += reset
		spriteLines = append(spriteLines, line)
	}
	return spriteLines, nil
}

func addCaught(id int, name string) error {
	caughtFile := filepath.Join(cacheDir, "caught.json")

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	caught := make(map[int]string)
	if data, err := os.ReadFile(caughtFile); err == nil {
		_ = json.Unmarshal(data, &caught)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read caught index: %w", err)
	}

	if _, exists := caught[id]; exists {
		return nil
	}

	caught[id] = name

	keys := make([]int, 0, len(caught))
	for k := range caught {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	sorted := make(map[string]string)
	for _, k := range keys {
		sorted[fmt.Sprintf("%d", k)] = caught[k]
	}

	out, err := json.MarshalIndent(sorted, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal caught file: %w", err)
	}

	return os.WriteFile(caughtFile, out, 0644)
}

func IsCaught(name string) (bool, error) {
	caughtFile := filepath.Join(cacheDir, "caught.json")
	data, err := os.ReadFile(caughtFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to read file: %w", err)
	}

	caught := make(map[string]string)
	if err := json.Unmarshal(data, &caught); err != nil {
		return false, fmt.Errorf("failed to unmarshal caught.json: %w", err)
	}

	for _, caughtName := range caught {
		if caughtName == name {
			return true, nil
		}
	}
	return false, nil
}
