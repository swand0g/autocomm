package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

type model struct {
	keymap   					keymap
	
	commitMsgCursor int
	selected 				map[int]struct{}
	commitChoices  	[]string
	
	aiModelCursor		int
	aiModel					string
	aiModels 			 	[]string
	maxTokens       int
	useConventional bool

	apiKey        string
	authenticated bool

	fetching   				 			 bool
	fetchError 				 			 bool
	shouldRefetchForNewModel bool

	appstate    int
	commitState commitState

	help     	help.Model
	spinner   spinner.Model
	textInput textinput.Model
}

type commitState struct {
	chosenMsg    string
	committed    bool
	committing   bool
	commitOutput string
	err          error
}

type keymap struct {
	Quit         	key.Binding
	Enter        	key.Binding
	Up           	key.Binding
	Down         	key.Binding
	Escape       	key.Binding
	Retry        	key.Binding
	Authenticate 	key.Binding
	ChooseAIModel key.Binding
}

type (
	requestStrResponse    struct{ data string }
	requestResponse struct{ data []string }
	requestError          struct{ err error }

	commitResult struct {
		output string
		err    error
	}
)

const (
	Choosing       = iota
	Authenticating = iota
	Quitting       = iota
	Committing     = iota
	ChoosingAIModel  = iota
)
