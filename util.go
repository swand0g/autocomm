package main

import (
	"time"
	"log"
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

func gitDiff() string {
	out, err := exec.Command("git", "diff").Output()
	if err != nil {
			log.Fatal(err)
			panic("some error found")
	}
	return string(out)
}