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

	//kick off chat service
	go chat.SetupWebService()
	session.SetupSession()
}
