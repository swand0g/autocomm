package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type env struct {
	DEBUG bool
	DRY 	bool
}

var environment env

// todo: fix logging when in build mode
func setupLogging() *os.File {
	if environment.DEBUG {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		return f
	}
	return nil
}

// todo: handle exit code error 129 when run in a folder that isnt a git repo
func main() {
	environment = env{
		DEBUG: len(os.Getenv("DEBUG")) > 0,
		DRY: len(os.Getenv("DRY")) > 0,
	}

	f := setupLogging()
	if f != nil { defer f.Close() }

	prog := tea.NewProgram(InitalModel())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
