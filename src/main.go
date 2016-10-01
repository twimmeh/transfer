package main

import (
	"flag"
	"session"
)

func main() {
	flag.Parse()
	session.SetupSession()
}
