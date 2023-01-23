package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
				case key.Matches(msg, m.keymap.Quit):
					m.appstate = Quitting
					return m, tea.Quit

				case key.Matches(msg, m.keymap.Up):
					if m.cursor > 0 {
						m.cursor--
					}

				case key.Matches(msg, m.keymap.Down):
					if m.cursor < len(m.choices) - 1 {
						m.cursor++
					}

				case key.Matches(msg, m.keymap.Choose):
					_, ok := m.selected[m.cursor]
					if ok {
						delete(m.selected, m.cursor)
					} else {
						m.selected[m.cursor] = struct{}{}
					}
				
				case key.Matches(msg, m.keymap.Authenticate):
					m.appstate = Authenticating
			}
			

		default:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.appstate {
		case Quitting:
			todo()

		case Authenticating:
			todo()

		default:
			todo()
	}

	s := "Choose a commit message...\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">" //cursor!
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	
	s += m.HelpView()
	if (m.appstate == Authenticating) {
		peaceOutMsg := randGoodbyeMessage()
		return fmt.Sprintf("\n%s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(peaceOutMsg))
	}

	return s
}



func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func InitalModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
	token, err := getAuthToken()
	authenticated := token != "" && err != nil

	appstate := Authenticating
	if authenticated {
		appstate = Choosing
	}

	return model{
		token: token,
		authenticated: authenticated,
		choices: []string{"carrots", "celery", "beans"},
		selected: make(map[int]struct{}),
		keymap: keymap{
			Quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c", "esc"),
				key.WithHelp("q", "quit"),
			),
			Choose: key.NewBinding(
				key.WithKeys("enter", " "),
				key.WithHelp("enter", "select"),
			),
			Up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "down"),
			),
			Authenticate: key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "authenticate"),
			),
		},
		// views & components
		spinner: s,
		help: help.NewModel(),
		// app state
		appstate: appstate,
	}
}