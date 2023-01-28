package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var netsema = make(chan struct{}, 1)

func main() {
	os.Setenv("DEBUG", "1")
	f := setupLogging()
	if f != nil {
		defer f.Close()
	}

	p := tea.NewProgram(InitalModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func setupLogging() *os.File {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		return f
	}
	return nil
}

func printCommitSugestionsAs() {
	k, _ := getAPIKey()
	
	s, e := fetchCommitSuggestions(k, false)
	if e != nil {
		fmt.Println(e)
	}
	
	fmt.Println(s)
}
