package main

func resetCommitSuggestions(m *model) {
	m.commitChoices = []string{}
	m.fetchError = false
}