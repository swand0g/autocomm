package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)


func main() {
	p := tea.NewProgram(InitalModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func main2() {
	k, _ := getAPIKey()
	
	s, e := fetchCommitSuggestions(k, false)
	if e != nil {
		fmt.Println(e)
	}
	
	fmt.Println(s)
}
