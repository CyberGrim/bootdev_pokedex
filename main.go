package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/cybergrim/bootdev_pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	nextURL *string
	prevURL *string
	api     *pokeapi.Client
	pokedex map[string]pokeapi.PokemonInfo
}

var callback_registry map[string]cliCommand

func main() {
	initialURL := "https://pokeapi.co/api/v2/location-area/"
	cfg := &config{
		nextURL: &initialURL,
		api:     pokeapi.NewClient(5 * time.Second),
		pokedex: make(map[string]pokeapi.PokemonInfo),
	}
	callback_registry = map[string]cliCommand{
		"map":     {name: "map", description: "Go forward on the Pokemon World Map - 20 results", callback: commandMap},
		"mapb":    {name: "mapb", description: "Go backwards on the Pokemon World Map - 20 results", callback: commandMapb},
		"catch":   {name: "catch", description: "Catch a particular Pokemon - eg. `catch <Pokemon>`", callback: commandCatch},
		"inspect": {name: "inspect", description: "Inspect a particular Pokemon - eg. `inspect <Pokemon>`", callback: commandInspect},
		"explore": {name: "explore", description: "Explore a particular area - eg. `explore <area>`", callback: commandExplore},
		"exit":    {name: "exit", description: "Exit the Pokedex", callback: commandExit},
		"help":    {name: "help", description: "Displays a help message", callback: commandHelp},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() == false {
			break
		}
		input := scanner.Text()

		input_cleaned := strings.TrimSpace(input)
		input_cleaned = strings.ToLower(input_cleaned)

		word_list := strings.Fields(input_cleaned)

		if len(word_list) == 0 {
			continue
		} else {
			command := word_list[0]
			area := word_list[1:]
			req_callback, ok := callback_registry[command]
			if !ok {
				fmt.Println("Unknown command")
				continue
			} else {
				err := req_callback.callback(cfg, area)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func commandMap(cfg *config, area []string) error {
	if cfg.nextURL == nil {
		return fmt.Errorf("you are on the last page")
	}
	return fetchAndPrintLocations(cfg, *cfg.nextURL)
}

func commandMapb(cfg *config, area []string) error {
	if cfg.prevURL == nil {
		return fmt.Errorf("you are on the first page")
	}
	return fetchAndPrintLocations(cfg, *cfg.prevURL)
}

func fetchAndPrintLocations(cfg *config, url string) error {
	res, err := cfg.api.ListLocationAreas(&url)
	if err != nil {
		return err
	}
	cfg.nextURL = res.Next
	cfg.prevURL = res.Prev

	for _, area := range res.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandCatch(cfg *config, pokemon []string) error {
	if len(pokemon) == 0 {
		return fmt.Errorf("you must provide the name of a Pokemon")
	}
	pokemonName := pokemon[0]
	res, err := cfg.api.GetPokemonInfo(pokemonName)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	catchChance := rand.Intn(res.BaseExperience)
	catchThreshhold := 40
	if catchChance >= catchThreshhold {
		fmt.Printf("%s escaped!\n", pokemonName)
	} else {
		fmt.Printf("%s was caught!\n", pokemonName)
		cfg.pokedex[pokemonName] = res
	}

	return nil
}

func commandInspect(cfg *config, pokemon []string) error {
	pokemonName := pokemon[0]
	entry, exists := cfg.pokedex[pokemonName]
	if exists {
		fmt.Printf("Name: %s\n", entry.Name)
		fmt.Printf("Height: %d\n", entry.Height)
		fmt.Printf("Weight: %d\n", entry.Weight)
		fmt.Println("Stats:")
		for _, stat := range entry.Stats {
			fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, ret_types := range entry.Types {
			fmt.Printf("  -%s\n", ret_types.Type.Name)
		}
	} else {
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}

func commandExplore(cfg *config, area []string) error {
	if len(area) == 0 {
		return fmt.Errorf("you must provide a location name")
	}
	locationName := area[0]
	res, err := cfg.api.ExploreLocation(locationName)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", locationName)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range res.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandExit(cfg *config, area []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, area []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	keys := []string{"map", "mapb", "catch", "inspect", "explore", "exit", "help"}
	for _, k := range keys {
		cmd := callback_registry[k]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}
