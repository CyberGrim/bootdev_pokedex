package main

import "strings"

func cleanInput(text string) []string {
	lowered := strings.ToLower(strings.TrimSpace(text))
	return strings.Fields(lowered)
}
