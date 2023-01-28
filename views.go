package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
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
				m.keymap.Enter,
				m.keymap.Authenticate,
			)
			break
		case Authenticating:
			esc := m.keymap.Escape
			esc.SetHelp(HelpText("esc"), "go back")

			keys = []key.Binding{
				esc,
			}
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
	if len(m.choices) == 0 && !m.fetchError {
		fs := fmt.Sprintf("\n%s %s\n", m.spinner.View(), "Fetching commit messages...")
		return fs
	} else if m.fetchError {
		es := fmt.Sprintf("%s", "Error fetching commit messages! Check your API key and try again.")
		return es 
	}

	s := "\n"
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
