package main

import (
	"fmt"
	"log"
	"math/rand"
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

func cleanLines(str string) []string {
	lines := strings.Split(str, "\n")
	re := regexp.MustCompile(`^(\d+\.|-|\*)\s+`)
	for i, line := range lines {
		lines[i] = strings.TrimSpace(re.ReplaceAllString(line, ""))
	}
	return lines
}

func userInGitRepo() bool {
	_, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	return err == nil
}
