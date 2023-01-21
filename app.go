package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices 	[]string
	cursor 		int
	selected 	map[int]struct{}
	help 		  help.Model
	keymap 		keymap

	quitting  bool

	spinner  spinner.Model
}

type keymap struct {
	Quit 		key.Binding
	Choose  key.Binding
	Up 			key.Binding
	Down 		key.Binding
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			m.quitting = true
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
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	s := "What should we buy at the market?\n\n"
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

	
	s += fmt.Sprintf("\nPress q to quit! %s\n", m.spinner.View())
	s += m.helpView()
	
	if (m.quitting) {
		peaceOutMsg := randGoodbyeMessage()
		return s + fmt.Sprintf("\n\n%s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(peaceOutMsg))
	}

	return s
}


func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.Quit,
		m.keymap.Up,
		m.keymap.Down,
		m.keymap.Choose,
	})
}

func InitalModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		choices: []string{"carrots", "celery", "beans"},
		selected: make(map[int]struct{}),
		spinner: s,
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
		},
		help: help.NewModel(),
	}
}