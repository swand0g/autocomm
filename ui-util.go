package main

import (
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

func HelpText(s string) string {
	return lipgloss.NewStyle().
		Italic(true).
		Bold(true).
		Render(s)
}

func (m model) quitApp() (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m model) getCommitSuggestions() {
	data, err := fetchCommitSuggestions(m.apiKey, m.useConventional)
	if err != nil {
		m.fetchError = true
	}

	m.choices = data
	m.fetching = false
	m.fetchError = false
}