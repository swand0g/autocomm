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
			m.keymap.Retry,
			m.keymap.Authenticate,
			m.keymap.ChooseAIModel,
		)
		break
	case Authenticating:
		esc := m.keymap.Escape
		esc.SetHelp(HelpText("esc"), "go back")

		enter := m.keymap.Enter
		enter.SetHelp(HelpText("enter"), "authenticate")

		keys = []key.Binding{
			esc,
			enter,
		}
		break
	case ChoosingAIModel:
		esc := m.keymap.Escape
		esc.SetHelp(HelpText("esc"), "go back")

		enter := m.keymap.Enter
		enter.SetHelp(HelpText("enter"), "select")

		keys = []key.Binding{
			esc,
			enter,
		}
		break
	default:
		break
	}

	return "\n" + m.help.ShortHelpView(keys)
}

// todo: allow this to be configurable
func (m model) QuitView() string {
	peaceOutMsg := randGoodbyeMessage()
	return fmt.Sprintf("\n%s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Purple)).Render(peaceOutMsg))
}

func (m model) AuthenticatingView() string {
	prompt := "What's your OpenAI API key?"
	ti := m.textInput.View()
	return fmt.Sprintf("\n%s %s  %s\n", m.spinner.View(), prompt, ti)
}

func (m model) ChooseView() string {
	if len(m.commitChoices) == 0 && !m.fetchError && m.fetching {
		fs := fmt.Sprintf("  %s %s", m.spinner.View(), "Fetching commit messages...")
		return "\n" + fs + "\n"
	} else if m.fetchError {
		es := fmt.Sprintf(
			"  %s %s",
			TextWithColor("Error fetching commit messages!", colors.Red),
			"Check your API key and try again.",
		)
		return "\n" + es + "\n"
	}

	s := "\n"
	for i, choice := range m.commitChoices {
		cursor := " "
		if m.commitMsgCursor == i {
			// todo: implement view padding
			cursor = TextWithColor("  >", colors.Purple) //cursor!
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return s
}

func (m model) CommitView() string {
	if m.commitState.err != nil {
		return fmt.Sprintf("\n%s %s\n\n", TextWithColor("Error committing!", colors.Red), m.commitState.err)
	}

	if m.commitState.committed {
		return m.commitState.commitOutput
	}

	return fmt.Sprintf("\n%s %s\n\n", m.spinner.View(), "Committing...")
}

func (m model) ChooseAIModelView() string {
	s := "\n"

	for i, model := range m.aiModels {
		cursor := " "
		if m.aiModelCursor == i {
			cursor = TextWithColor("  >", colors.Purple)
		}

		mt := model
		if model == m.aiModel {
			mt = TextWithColor(model, colors.Tron)
		}
		s += fmt.Sprintf("%s %s\n", cursor, mt)
	}

	return s
}
