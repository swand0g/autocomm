package main

const CONFIG_FILE_NAME = ".autocomm"

var DRY_COMMIT_SUGGESTIONS = []string{
	"feat: add 'comments' option",
	"docs: remove reference to 'timeTravel'",
	"refactor: share logic between 4d3d3d3 and flarhgunnstow",
	"chore: release patch version 1.0.1",
	"revert: remove broken confirmation message",
}


var API_MODELS = []string{
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
