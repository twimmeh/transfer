package main

import (
	"chat"
	"config"
	"flag"
	"session"
)

func main() {

	flag.Parse()
	config.Parse()

	//kick off chat loop
	go chat.SendChatMessageLoop()
	session.SetupSession()
}
