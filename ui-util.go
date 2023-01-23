package main

import (
	"github.com/charmbracelet/lipgloss"
)

func HelpText(s string) string {
	return lipgloss.NewStyle().
		Italic(true).
		Bold(true).
		Render(s)
}