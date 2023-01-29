package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

type model struct {
	choices 	[]string
	cursor 		int
	selected 	map[int]struct{}
	help 		  help.Model
	keymap 		keymap

	apiKey 				string
	authenticated bool

	fetching 	 			bool
	fetchError 			bool
	maxTokens 			int
	useConventional bool

	appstate 		int
	commitState commitState

	spinner  	spinner.Model
	textInput textinput.Model
}

type commitState struct {
	chosenMsg  	 string
	committed  	 bool
	committing 	 bool
	commitOutput string
	err        	 error
}

type keymap struct {
	Quit 						key.Binding
	Enter  					key.Binding
	Up 							key.Binding
	Down 						key.Binding
	Authenticate 		key.Binding
	Escape 					key.Binding
	Retry 					key.Binding
}

type (
	requestStrResponse 		struct{ data string }
	requestStrArrResponse struct{ data []string }
	requestError		 	 		struct{ err error }

	commitResult struct{
		output string
		err		 error
	}
)

const (
	Choosing 				= iota
	Authenticating 	= iota
	Quitting 				= iota
	Committing 			= iota
)