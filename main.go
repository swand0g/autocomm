package main

import (
	"os"
	"fmt"
	"context"

	tea 	"github.com/charmbracelet/bubbletea"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func main() {
	// fmt.Println(gitDiff())
	fileName := ".autocomm"
	// writeToFile("hello world!!!", fileName)

	openaiToken, err := readFromFile(fileName)
	if err != nil { panic(err) }
	fmt.Println(openaiToken)
}

func main3() {
	// get the token from file
	token, err := readFromFile(".autocomm")
	if err != nil { panic(err) }

	c := gogpt.NewClient(token)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model: gogpt.GPT3TextDavinci003,
		MaxTokens: 25,
		Prompt:
		`The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.

		Human: Hello, who are you?
		AI: I am an AI created by OpenAI. How can I help you today?
		Human:`,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		println(err.Error())
		return
	}

	fmt.Println(resp.Choices[0].Text)
}

func main2() {
	p := tea.NewProgram(InitalModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}