package main

import (
	"context"

	gogpt "github.com/sashabaranov/go-gpt3"
)

const prompt = `
	The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.
	Human: Hello, who are you?
	AI: I am an AI created by OpenAI. How can I help you today?
	Human:
`

func fetchAiResponse(prompt string) (string, error) {
	token, err := readFromFile(ConfigFileName)
	if err != nil { return "", err }

	c := gogpt.NewClient(token)
	ctx := context.Background()
	req := gogpt.CompletionRequest{
		Model: gogpt.GPT3TextDavinci003,
		MaxTokens: 25,
		Prompt: prompt,
	}
	
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil { return "", err }

	return resp.Choices[0].Text, nil
}