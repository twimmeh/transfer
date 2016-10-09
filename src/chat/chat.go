package chat

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"session"
	"strings"
)

const chatServiceId = 1

//Register chat service
func init() {
	chatService := ChatService{}
	session.RegisterService(chatServiceId, chatService)
}

type ChatService struct{}

// Simple read and dump to command line
func (svc ChatService) HandleConnection(session session.Session, conn session.Connection) {
	b, _ := ioutil.ReadAll(conn)
	fmt.Println(string(b[:]))
}

// Poll console for messages to send
func SendChatMessageLoop() {
	for {
		sendMessage(getChatMessage())
	}
}

// If a session is available, try to send a message to it
func sendMessage(message string) {
	s, _ := session.GetSession()
	if s != nil {
		conn := s.OpenConnection(chatServiceId)
		io.WriteString(conn, message)
		conn.Close()
	} else {
		fmt.Println("No chat sessions available")
	}
}

// Helper function to retreive text message from console
func getChatMessage() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
