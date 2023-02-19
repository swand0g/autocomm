package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func saveConfig(token string) bool {
	err := viper.WriteConfig()
	if err != nil {
		logi("Error saving config: %v", err)
	}
	return err == nil
}

func getAPIKey() (string, error) {
	key := viper.GetString("apiKey")
	if key == "" {
		return "", fmt.Errorf("No api key found")
	}
	return key, nil
}

func getAIModel() (string, error) {
	model := viper.GetString("aiModel")
	if model == "" {
		return "", fmt.Errorf("No ai model found")
	}
	return model, nil
}
