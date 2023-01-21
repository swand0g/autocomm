package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
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
