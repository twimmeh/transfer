package main

import (
	"chat"
	"config"
	"flag"
	"session"
	"net/http"
)

func main() {

	flag.Parse()
	config.Parse()

	//kick off chat service
	chatPort := chat.SetupWebService()
	go http.ListenAndServe(":"+chatPort, nil)
	session.SetupSession()
}
