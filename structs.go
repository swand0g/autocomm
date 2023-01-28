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
	channel	    chan asyncMsg
	count int

	apiKey 				string
	authenticated bool

	appstate int
	
	fetching 	 			bool
	fetchError 			bool
	maxTokens 			int
	useConventional bool

	spinner  	spinner.Model
	textInput textinput.Model

	fet bool
}

type keymap struct {
	Quit 						key.Binding
	Enter  					key.Binding
	Up 							key.Binding
	Down 						key.Binding
	Authenticate 		key.Binding
	Escape 					key.Binding
}

type asyncMsg struct {
	choices []string
	count   int
}

type (
	requestStrResponse 		struct{ data string }
	requestStrArrResponse struct{ data []string }
	requestError		 	 		struct{ err error }
)

const (
	Choosing 				= iota
	Authenticating 	= iota
	Quitting 				= iota
)