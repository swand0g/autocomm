package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/key"
)

func (m model) HelpView() string {
	keys := []key.Binding{
		m.keymap.Quit,
	}

	switch m.appstate {
		case Choosing:
			keys = append(keys,
				m.keymap.Up,
				m.keymap.Down,
				m.keymap.Choose,
				m.keymap.Authenticate,
			)
			break
		default:
			break
	}
	
	return "\n" + m.help.ShortHelpView(keys)
}

func (m model) QuitView() string {
	peaceOutMsg := randGoodbyeMessage()
	return fmt.Sprintf("\n%s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(peaceOutMsg))
}

func (m model) AuthenticatingView() string {
	prompt := "What's your OpenAI API key?"
	ti := m.textInput.View()
	return fmt.Sprintf("\n%s %s  %s\n", m.spinner.View(), prompt, ti)
}

func (m model) ChooseView() string {
	s := "Choose a commit message...\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">" //cursor!
			// cursor = m.spinner.View()
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	return s
}