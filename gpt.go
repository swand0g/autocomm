package main

import (
	"context"
	"fmt"

	gogpt "github.com/sashabaranov/go-gpt3"
)

const PROMPT = "Suggest %d good commit messages for my commit%s:" +
	"```\n" + "%s" + "```\n"

func fetchAiResponse(apiKey string, prompt string, aiModel string) (string, error) {
	c := gogpt.NewClient(apiKey)
	ctx := context.Background()
	req := gogpt.CompletionRequest{
		Model:     aiModel,
		MaxTokens: 100,
		Prompt:    prompt,
	}

	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Text, nil
}

func fetchCommitSuggestions(apiKey string, conventional bool, aiModel string) ([]string, error) {
	diff, err := gitDiff()
	if err != nil {
		return []string{}, err
	}

	cs := ""
	if conventional {
		cs = " in conventional commit format (type(scope): subject)"
	}

	fullPrompt := fmt.Sprintf(PROMPT, 5, cs, diff)
	res, err := fetchAiResponse(apiKey, fullPrompt, aiModel)
	if err != nil {
		return []string{}, err
	}

	s := cleanLines(res)
	return s, nil
}
