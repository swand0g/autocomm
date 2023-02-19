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
	"github.com/spf13/viper"

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
	return requestResponse{[]string{r}}
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
		case requestResponse:
			logi("received res in Update(msg): %v", msg.data)
			m.fetching = false
			if !m.shouldRefetchForNewModel {
				m.commitChoices = msg.data
			}
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

		case ChoosingAIModel: {
			switch msg := msg.(type) {
				case tea.KeyMsg:
					switch {
						case key.Matches(msg, m.keymap.Quit, m.keymap.Escape):
							m.appstate = Choosing
							break
						
						case key.Matches(msg, m.keymap.Up):
							if m.aiModelCursor > 0 {
								m.aiModelCursor--
							} else {
								m.aiModelCursor = len(m.aiModels) - 1
							}
							break
						
						case key.Matches(msg, m.keymap.Down):
							if m.aiModelCursor < len(m.aiModels) - 1 {
								m.aiModelCursor++
							} else {
								m.aiModelCursor = 0
							}
							break
						
						case key.Matches(msg, m.keymap.Enter):
							m.aiModel = m.aiModels[m.aiModelCursor]
							m.shouldRefetchForNewModel = true
							m.appstate = Choosing
							logi("selected ai model: %v", m.aiModel)
							break

					}
			}

			break
		}
		
		case Choosing: {
			canFetch := !m.fetching && !m.fetchError
			shouldRefetch := len(m.commitChoices) == 0 || m.shouldRefetchForNewModel
			
			if shouldRefetch && canFetch {
				m.fetching = true
				if m.shouldRefetchForNewModel {
					m.shouldRefetchForNewModel = false
				}

				// used for debugging
				// cmds = append(cmds, req)
				
				cmds = append(cmds, m.getCommitSuggestions)
			}

			switch msg := msg.(type) {
				case tea.KeyMsg:
					switch {
						case key.Matches(msg, m.keymap.Quit, m.keymap.Escape):
							m.appstate = Quitting
							return m, tea.Quit

						case key.Matches(msg, m.keymap.Up):
							if m.commitMsgCursor > 0 {
								m.commitMsgCursor--
							} else {
								m.commitMsgCursor = len(m.commitChoices) - 1
							}
							break

						case key.Matches(msg, m.keymap.Down):
							if m.commitMsgCursor < len(m.commitChoices) - 1 {
								m.commitMsgCursor++
							} else {
								m.commitMsgCursor = 0
							}
							break

						case key.Matches(msg, m.keymap.Retry):
							if !m.fetching {
								m.fetchError = false
								m.commitChoices = []string{}
							}
							break

						case key.Matches(msg, m.keymap.Enter):
							logi("selected: %v", m.commitChoices[m.commitMsgCursor])
							_, ok := m.selected[m.commitMsgCursor]
							if ok {
								delete(m.selected, m.commitMsgCursor)
							} else {
								m.selected[m.commitMsgCursor] = struct{}{}
							}

							cmtMsg := m.commitChoices[m.commitMsgCursor]
							m.commitState.committing = true
							m.commitState.chosenMsg = cmtMsg
							m.appstate = Committing
							cmds = append(cmds, m.commitWithMsg)
							break

						case key.Matches(msg, m.keymap.Authenticate):
							m.appstate = Authenticating
							break
						
						case key.Matches(msg, m.keymap.ChooseAIModel):
							m.appstate = ChoosingAIModel
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
							viper.Set("apiKey", key)
							saveConfig(key)
							m.apiKey = key
							m.fetchError = false
							m.commitChoices = []string{}
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
		
		case ChoosingAIModel:
			v += m.ChooseAIModelView()
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
	initiallyAuthenticated := apiKey != "" && err == nil
	
	appstate := Authenticating
	if initiallyAuthenticated {
		appstate = Choosing
	}

	models := []string{
		"text-davinci-003",
		"text-davinci-002",
		"text-davinci-001",
		"text-curie-001",
		"text-babbage-001",
		"text-ada-001",
		"code-davinci-002",
		"code-cushman-001",
		"code-davinci-001",
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
		authenticated: initiallyAuthenticated,
		commitChoices: []string{},
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
				key.WithKeys("enter"),
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
			Retry: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp(HelpText("r"), "retry"),
			),
			Authenticate: key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp(HelpText("a"), "set auth"),
			),
			ChooseAIModel: key.NewBinding(
				key.WithKeys("m"),
				key.WithHelp(HelpText("m"), "ai model"),
			),
		},
		// data
		aiModel: "text-davinci-003",
		aiModels: models,
		// views & components
		spinner: s,
		help: h,
		textInput: ti,
		// app state
		appstate: appstate,
		fetching: false,
		fetchError: false,
		shouldRefetchForNewModel: false,
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
