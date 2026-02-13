package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var callback_registry map[string]cliCommand

func main() {
	callback_registry = map[string]cliCommand{
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
				err := req_callback.callback()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
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
