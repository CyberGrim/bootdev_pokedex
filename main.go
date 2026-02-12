package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
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
			fmt.Printf("Your command was: %s\n", command)
		}
	}
}
