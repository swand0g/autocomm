package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func req() tea.Msg {
	time.Sleep(1 * time.Second)
	url := "https://jsonplaceholder.typicode.com/todos/1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return requestError{err}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return requestError{err}
	}

	r := strings.Replace(string(body), "\n", "", -1)
	log.Println("fetched: ", r)
	return requestStrArrResponse{[]string{r}}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds = []tea.Cmd{}
		spinCmd tea.Cmd
	)

	m.spinner, spinCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinCmd)

	// Handle async messages first
	switch msg := msg.(type) {
		case requestStrArrResponse:
			logi("GOT RES IN UPDATE(msg): %v", msg.data)
			m.choices = msg.data
			m.fetching = false
			break
		case requestError:
			m.fetchError = true
			m.fetching = false
			break
		case commitResult:
			if msg.err != nil {
				m.commitState.err = msg.err
			} else {
				m.commitState.committed = true
				m.commitState.committing = false
				m.commitState.commitOutput = msg.output
			}
			break
	}

	switch m.appstate {
		case Choosing: {
			if len(m.choices) == 0 && !m.fetching && !m.fetchError {
				m.fetching = true
				cmds = append(cmds, m.getCommitSuggestions)

				// used for debugging
				// cmds = append(cmds, req)
			}

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
							break

						case key.Matches(msg, m.keymap.Down):
							if m.cursor < len(m.choices) - 1 {
								m.cursor++
							} else {
								m.cursor = 0
							}
							break

						case key.Matches(msg, m.keymap.Retry):
							if !m.fetching {
								m.fetchError = false
								m.choices = []string{}
							}
							break

						case key.Matches(msg, m.keymap.Enter):
							logi("selected: %v", m.choices[m.cursor])
							_, ok := m.selected[m.cursor]
							if ok {
								delete(m.selected, m.cursor)
							} else {
								m.selected[m.cursor] = struct{}{}
							}

							cmtMsg := m.choices[m.cursor]
							m.commitState.committing = true
							m.commitState.chosenMsg = cmtMsg
							m.appstate = Committing
							cmds = append(cmds, m.commitWithMsg)
							break

						case key.Matches(msg, m.keymap.Authenticate):
							m.appstate = Authenticating
							break
					}

				default:
					break
			}
		}

		case Authenticating: {
			var textInputCmd tea.Cmd

			switch msg := msg.(type) {
				case tea.KeyMsg:
					switch {
						case key.Matches(msg, m.keymap.Escape):
							m.textInput.Reset()
							m.appstate = Choosing
							break
						case key.Matches(msg, m.keymap.Enter):
							key := m.textInput.Value()
							saveAPIKey(key)
							m.apiKey = key
							m.fetchError = false
							m.choices = []string{}
							m.textInput.Reset()
							m.appstate = Choosing
							break
					}
			}

			m.textInput, textInputCmd = m.textInput.Update(msg)
			cmds = append(cmds, textInputCmd)

			m.textInput.Focus()
			break
		}

		case Committing:
			if m.commitState.err != nil || m.commitState.committed {
				return m, tea.Quit
			}
			break

		default:
			break
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	v := ""

	switch m.appstate {
		case Quitting:
			return m.QuitView()

		case Authenticating:
			v += m.AuthenticatingView()
			break

		case Committing:
			v += m.CommitView()
			if (m.commitState.committed) {
				return v
			}
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

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
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

	ti := textinput.New()
	ti.Placeholder = "Super secret API key"
	ti.CharLimit = 64
	ti.Width = 64

	h := help.NewModel()

	return model{
		apiKey: apiKey,
		authenticated: authenticated,
		choices: []string{},
		selected: make(map[int]struct{}),
		maxTokens: 100,
		useConventional: false,
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
			Retry: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp(HelpText("r"), "retry"),
			),
		},
		// views & components
		spinner: s,
		help: h,
		textInput: ti,
		// app state
		appstate: appstate,
		fetching: false,
		fetchError: false,
		// substates
		commitState: commitState{
			chosenMsg: "",
			committed: false,
			committing: false,
			commitOutput: "",
			err: nil,
		},
	}
}
