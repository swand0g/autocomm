package main

import (
	"fmt"
	"time"

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

// func listenForActivity(sub chan struct{}) tea.Cmd {
// 	return func() tea.Msg {
// 		for {
// 			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(900)+100))
// 			sub <- struct{}{}
// 		}
// 	}
// }

func (m model) f() {
	time.Sleep(time.Second * 1)
	m.channel <- asyncMsg{ count: m.count + 1 }
}

func (m model) ChooseView() string {
	if len(m.choices) == 0 && !m.fetchError {
		if !m.fetching {
			// todo: rip my grant lmao this ran like a million requests
			// go m.getCommitSuggestions()
			// m.fetching = true

			// m.fetching = true
		}

		if m.fetching {
			go m.f()
		}

		fs := fmt.Sprintf("\n%s %s: %d\n", m.spinner.View(), "Fetching commit messages...", m.count)
		return fs
	} else if m.fetchError {
		es := fmt.Sprintf("%s", "Error fetching commit messages! Check your API key and try again.")
		return es 
	}

	s := ""
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
