package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type env struct {
	DEBUG            bool
	DRY              bool
	USE_CONVENTIONAL bool
}

var environment env

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

func main() {
	if !userInGitRepo() {
		fmt.Println("This ain't a git repo 🤨")
		os.Exit(1)
	}

	environment = env{
		DEBUG: false,
		DRY: false,
		USE_CONVENTIONAL: false,
	}

	flag.BoolVar(&environment.DEBUG, "debug", false, "enable debug logging")
	flag.BoolVar(&environment.DRY, "dry", false, "dry run (doesn't make API calls)")
	flag.BoolVar(&environment.USE_CONVENTIONAL, "conventional", false, "use conventional commits")
	flag.Parse()
	
	f := setupLogging()
	if f != nil { defer f.Close() }

	prog := tea.NewProgram(InitalModel())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
