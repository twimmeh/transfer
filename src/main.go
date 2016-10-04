package main

import (
	"chat"
	"flag"
	"session"
)

func main() {

	flag.Parse()

	//kick off chat loop
	go chat.SendChatMessageLoop()

	session.SetupSession()
}
