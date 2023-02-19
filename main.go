package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

var VERSION = "1.1"

type env struct {
	DEBUG            bool
	DRY              bool
	USE_CONVENTIONAL bool
}

var environment env

func setupLogging() *os.File {
	if environment.DEBUG {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		return f
	}
	return nil
}

func main() {
	if !userInGitRepo() {
		fmt.Println("This ain't a git repo 🤨")
		os.Exit(1)
	}

	environment = env{
		DEBUG: false,
		DRY: false,
		USE_CONVENTIONAL: false,
	}

	flag.BoolVar(&environment.DEBUG, "debug", false, "enable debug logging")
	flag.BoolVar(&environment.DRY, "dry", false, "dry run (doesn't make API calls)")
	flag.BoolVar(&environment.USE_CONVENTIONAL, "conventional", false, "use conventional commits")
	flag.Parse()
	
	setupConfig()
	f := setupLogging()
	if f != nil { defer f.Close() }

	prog := tea.NewProgram(InitalModel())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func createConfigFile() (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	dirPath := filepath.Join(home, ".config", "autocomm")
	mkdirErr := os.MkdirAll(dirPath, os.ModePerm)
	if mkdirErr != nil {
		return false, mkdirErr
	}

	filePath := filepath.Join(dirPath, "autocomm.toml")
	f, fileMakeErr := os.Create(filePath)
	defer f.Close()
	if err != nil {
		return false, fileMakeErr
	}

	return true, nil
}

func setupConfig() {
	viper.SetConfigName("autocomm.toml")
	viper.SetConfigType("toml")
	
	home, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil { fmt.Println("An unexpected error occurred when getting your home directory 😕") }

	viper.AddConfigPath(filepath.Join(home, ".config", "autocomm"))
	viper.AddConfigPath(home)

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if ok, createErr := createConfigFile(); !ok {
				fmt.Println("An unexpected error occurred when creating the config file 😕")
				logi(createErr.Error())
				os.Exit(1)
			}
			logi("Config file created at $HOME/.config/autocomm/.autocomm.toml")
		}
	} else {
		logi("Config file found!")
	}

	viper.SetDefault("apiKey", "")
	viper.SetDefault("aiModel", "text-davinci-003")
}

func main2() {
	setupConfig()
	fmt.Printf("%v\n", viper.AllSettings())
	// viper.WriteConfig()
}