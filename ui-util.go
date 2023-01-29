package main

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Color struct {
	Red     string
	Orange  string
	Yellow  string
	Green   string
	Blue    string
	Purple  string
	Pink    string
	Brown   string
	Black   string
	White   string
	Gray    string
	Silver  string
	Gold    string
}

var colors = Color{
	Red:     "#ff0000",
	Orange:  "#ff7f00",
	Yellow:  "#ffff00",
	Green:   "#00ff00",
	Blue:    "#0000ff",
	Purple:  "#8b00ff",
	Pink:    "#ff00ff",
	Brown:   "#a52a2a",
	Black:   "#000000",
	White:   "#ffffff",
	Gray:    "#808080",
	Silver:  "#c0c0c0",
	Gold:    "#ffd700",
}

/* Components */
func HelpText(s string) string {
	return lipgloss.NewStyle().
		Italic(true).
		Bold(true).
		Render(s)
}

func TextWithColor(s string, color string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render(s)
}

/* Misc */
func (m model) getCommitSuggestions() tea.Msg {
	data, err := fetchCommitSuggestions(m.apiKey, m.useConventional)
	
	logi("res for getCommitSuggestions(): %v", struct{ data []string; err string; apikey string; useConventional bool }{
		data: data,
		err: err.Error(),
		apikey: m.apiKey,
		useConventional: m.useConventional,
	})

	if err != nil { return requestError{err} }
	return requestStrArrResponse{data}
}

func (m model) commitWithMsg() tea.Msg {
	cmd := exec.Command("git", "commit", "-m", m.commitState.chosenMsg)
	out, _ := cmd.CombinedOutput()
	return commitResult{string(out), nil}
}
