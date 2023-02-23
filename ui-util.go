package main

import (
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Color struct {
	Red    string
	Orange string
	Yellow string
	Green  string
	Blue   string
	Purple string
	Pink   string
	Brown  string
	Black  string
	White  string
	Gray   string
	Silver string
	Gold   string
	Tron   string
}

var colors = Color{
	Red:    "#ff0000",
	Orange: "#ff7f00",
	Yellow: "#ffff00",
	Green:  "#00ff00",
	Blue:   "#0000ff",
	Purple: "#8b00ff",
	Pink:   "#ff00ff",
	Brown:  "#a52a2a",
	Black:  "#000000",
	White:  "#ffffff",
	Gray:   "#808080",
	Silver: "#c0c0c0",
	Gold:   "#ffd700",
	Tron:   "#00ffff",
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
	if environment.DRY {
		time.Sleep(500 * time.Millisecond)
		return requestResponse{DRY_COMMIT_SUGGESTIONS}
	}

	data, err := fetchCommitSuggestions(m.apiKey, m.useConventional, m.aiModel)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	logi("res for getCommitSuggestions(): %v", struct {
		data            []string
		err             string
		apikey          string
		useConventional bool
	}{
		data:            data,
		err:             errStr,
		apikey:          m.apiKey,
		useConventional: m.useConventional,
	})

	goodCommitSuggestions := []string{}
	for _, suggestion := range data {
		if suggestion != "" {
			goodCommitSuggestions = append(goodCommitSuggestions, suggestion)
		}
	}

	if err != nil {
		return requestError{err}
	}
	return requestResponse{goodCommitSuggestions}
}

func (m model) commitWithMsg() tea.Msg {
	cmd := exec.Command("git", "commit", "-m", m.commitState.chosenMsg)
	out, _ := cmd.CombinedOutput()
	return commitResult{string(out), nil}
}
