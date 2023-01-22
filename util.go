package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"time"
	"os/exec"
	"math/rand"
)

func randGoodbyeMessage() string {
	rand.Seed(time.Now().UnixNano())
	messages := []string{
		"Bye!",
		"Peace out!",
		"See ya!",
		"Ciao!",
		"Adios!",
		"Later!",
		"Have a nice day!",
		"Have a good one!",
	}
	return messages[rand.Intn(len(messages))]
}

func gitDiff() (string, error) {
	out, err := exec.Command("git", "diff").Output()
	if err != nil {
			log.Fatal(err)
			return "", err
	}
	return string(out), nil
}

func writeToFile(text string, fileName string) bool {
	home, err := os.UserHomeDir()
	if err != nil { return false }

	filePath := fmt.Sprintf("%s/%s", home, fileName)
	f, err := os.Create(filePath)
	if err != nil { return false }
	
	defer f.Close()
	_, err = io.WriteString(f, text)
	return err != nil
}

func readFromFile(fileName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil { return "", err }

	filePath := fmt.Sprintf("%s/%s", home, fileName)
	f, err := os.Open(filePath)
	if err != nil { return "", err }
	defer f.Close()
	
	buf := make([]byte, 1024)
	_, err = f.Read(buf)
	return string(buf), err
}