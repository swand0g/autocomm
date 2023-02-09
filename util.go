package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func todo(params ...string) {
	if len(params) == 0 {
		log.Println("todo!")
	} else {
		log.Printf("todo! %s", strings.Join(params, ", "))
	}
}

func logi(format string, v ...interface{}) {
	if !environment.DEBUG { return }

	_, file, line, ok := runtime.Caller(1)
	fileInfo := strings.Split(file, "/")
	fileName := fileInfo[len(fileInfo)-1]

	if ok {
		log.Printf("[%s:%d] %s", fileName, line, fmt.Sprintf(format, v...))
	}
}

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
	out, err := exec.Command("git", "diff", "--cached").Output()
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
	buf, err := os.ReadFile(filePath)
	if err != nil { return "", err }

	s := strings.TrimSpace(string(buf))
	return s, err
}

func saveAPIKey(token string) bool {
	return writeToFile(token, CONFIG_FILE_NAME)
}

func getAPIKey() (string, error) {
	return readFromFile(CONFIG_FILE_NAME)
}

func cleanLines(str string) []string {
	lines := strings.Split(str, "\n")
	re := regexp.MustCompile(`^(\d+\.|-|\*)\s+`)
	for i, line := range lines {
			lines[i] = strings.TrimSpace(re.ReplaceAllString(line, ""))
	}
	return lines
}