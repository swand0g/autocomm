package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) debugRequest() {
	url := "https://jsonplaceholder.typicode.com/todos/1"
	resp, err := http.Get(url)
	if err != nil {
			fmt.Println(err)
			return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
			fmt.Println(err)
			return
	}

	fmt.Printf("\n\n%s\n\n", string(body))
	m.choices = []string{string(body)}
	m.fetching = false
}


func main() {
	p := tea.NewProgram(InitalModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func main3() {
	k, _ := getAPIKey()
	
	s, e := fetchCommitSuggestions(k, false)
	if e != nil {
		fmt.Println(e)
	}
	
	fmt.Println(s)
}
