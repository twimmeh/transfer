package main

import (
	"bufio"
	"flag"
	"os"
	"session"
	"strings"
	"chat"
)

func main() {

	flag.Parse()

	//kick off chat loop
	go sendChatMessageLoop()

	session.SetupSession()
}


func sendChatMessageLoop(){
	for{
		m:=getChatMessage()
		chat.SendMessage(m)
	}
}

//Get message from console
func getChatMessage() string{
	reader:= bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
