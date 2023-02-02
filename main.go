package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

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

func main() {
	f := setupLogging()
	if f != nil { defer f.Close() }

	prog := tea.NewProgram(InitalModel())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
