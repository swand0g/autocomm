package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var spinCmd tea.Cmd
	var cmds []tea.Cmd = []tea.Cmd{}

	// Handle async messages first
	switch msg := msg.(type) {
		case asyncMsg:
			m.count = msg.count
			// m.fetching = false
			cmds = append(cmds, receiveMessage(m.channel))
	}

	switch m.appstate {
		case Choosing: {
			switch msg := msg.(type) {
				case tea.KeyMsg:
					switch {
						case key.Matches(msg, m.keymap.Quit, m.keymap.Escape):
							m.appstate = Quitting
							return m, tea.Quit
		
						case key.Matches(msg, m.keymap.Up):
							if m.cursor > 0 {
								m.cursor--
							} else {
								m.cursor = len(m.choices) - 1
							}
		
						case key.Matches(msg, m.keymap.Down):
							if m.cursor < len(m.choices) - 1 {
								m.cursor++
							} else {
								m.cursor = 0
							}
		
						case key.Matches(msg, m.keymap.Enter):
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
					m.spinner, spinCmd = m.spinner.Update(msg)
					cmds = append(cmds, spinCmd)
			}
			
			return m, tea.Batch(cmds...)
		}
		
		case Authenticating: {
			switch msg := msg.(type) {
				case tea.KeyMsg:
					switch {
						case key.Matches(msg, m.keymap.Escape):
							m.appstate = Choosing
							break
						case key.Matches(msg, m.keymap.Enter):
							saveAPIKey(m.textInput.Value())
							m.appstate = Choosing
							break
					}
			}

			m.textInput, spinCmd = m.textInput.Update(msg)
			m.textInput.Focus()
			break
		}
		
		default:
			break
	}

	return m, nil
}

func (m model) View() string {
	v := ""

	switch m.appstate {
		case Quitting:
			return m.QuitView()

		case Authenticating:
			v += m.AuthenticatingView()
			break

		case Choosing:
			fallthrough
		default:
			v += m.ChooseView()
			break
	}
	
	v += m.HelpView()
	return v
}

func receiveMessage(c chan asyncMsg) tea.Cmd {
	return func() tea.Msg {
		return asyncMsg(<-c)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		receiveMessage(m.channel),
	)
}

func InitalModel() model {
	apiKey, err := getAPIKey()
	authenticated := apiKey != "" && err == nil
	
	appstate := Authenticating
	if authenticated {
		appstate = Choosing
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	h := help.NewModel()

	ti := textinput.New()
	ti.Placeholder = "Super secret API key"
	ti.CharLimit = 64
	ti.Width = 20

	return model{
		apiKey: apiKey,
		authenticated: authenticated,
		choices: []string{},
		selected: make(map[int]struct{}),
		maxTokens: 100,
		useConventional: false,
		channel: make(chan asyncMsg),
		keymap: keymap{
			Quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp(HelpText("q"), "quit"),
			),
			Escape: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp(HelpText("esc"), "escape"),
			),
			Enter: key.NewBinding(
				key.WithKeys("enter", " "),
				key.WithHelp(HelpText("enter"), "select"),
			),
			Up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp(HelpText("↑/k"), "up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp(HelpText("↓/j"), "down"),
			),
			Authenticate: key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp(HelpText("a"), "set auth"),
			),
		},
		// views & components
		spinner: s,
		help: h,
		textInput: ti,
		// app state
		appstate: appstate,
		fetching: true,
		fetchError: false,
	}
}
