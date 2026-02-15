package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cybergrim/bootdev_pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	nextURL *string
	prevURL *string
	api     *pokeapi.Client
}

var callback_registry map[string]cliCommand

func main() {
	initialURL := "https://pokeapi.co/api/v2/location-area/"
	cfg := &config{
		nextURL: &initialURL,
		api:     pokeapi.NewClient(5 * time.Second),
	}
	callback_registry = map[string]cliCommand{
		"map":  {name: "map", description: "Go forward on the Pokemon World Map - 20 results", callback: commandMap},
		"mapb": {name: "mapb", description: "Go backwards on the Pokemon World Map - 20 results", callback: commandMapb},
		"exit": {name: "exit", description: "Exit the Pokedex", callback: commandExit},
		"help": {name: "help", description: "Displays a help message", callback: commandHelp},
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
			req_callback, ok := callback_registry[command]
			if !ok {
				fmt.Println("Unknown command")
				continue
			} else {
				err := req_callback.callback(cfg)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func commandMap(cfg *config) error {
	if cfg.nextURL == nil {
		return fmt.Errorf("you are on the last page")
	}
	return fetchAndPrintLocations(cfg, *cfg.nextURL)
}

func commandMapb(cfg *config) error {
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

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	keys := []string{"help", "exit"}
	for _, k := range keys {
		cmd := callback_registry[k]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}
