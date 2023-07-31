package main

import (
	"log"

	twitch "github.com/rwxrob/cozy-twitch"
)

func main() {
	agent := &twitch.Agent{
		User:    `rwxbot`,
		Channel: `rwxrob`,
		Pass:    `oauth:ffesrf4v4rj11ryb8uhv506fojs648`,
	}
	err := agent.Connect()
	if err != nil {
		log.Fatal(err)
	}
	//	<-agent.Chan

}
