package main

import (
	"config"
	"flag"
	"session"
)

func main() {
	flag.Parse()
	config.Parse()
	session.SetupSession()
}
