package main

import (
	"github.com/charmbracelet/bubbles/key"
)

func (m model) HelpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.Quit,
		m.keymap.Up,
		m.keymap.Down,
		m.keymap.Choose,
	})
}