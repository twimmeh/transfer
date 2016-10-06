package chat

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"session"
	"strconv"
)

const chatServiceId = 1
const chatResource = "/chat"

var receivedChatMessages = ""
var chatPort = flag.String("chatport", "8080", "The HTTP port of the chat service. Default: 8080")

const chatPage = `<!DOCTYPE html>
<html>
<body>
<form>
  Chat message:<br>
  <input type="text" name="chatMessageInput">
  <br>
</form>
<div>
	<p>
		Other persons text to go here:
		<br>		
		%s
	</p>
</div>
</body>
</html>
`

//Register chat service
func init() {
	chatService := ChatService{}
	session.RegisterService(chatServiceId, chatService)
}

type ChatService struct{}

// Simple read and dump to command line
func (svc ChatService) HandleConnection(session session.Session, conn session.Connection) {

	b, _ := ioutil.ReadAll(conn)
	receivedChatMessages = receivedChatMessages + "\n" + string(b[:])
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

// Web-based chat message sending
func SetupWebService() {

	if _, err := strconv.Atoi(*chatPort); err != nil {
		panic("Invalid chat port specified (must be 32-bit integer). You put: " + *chatPort)
	}

	http.HandleFunc(chatResource, chatHandler)
	http.ListenAndServe(":"+*chatPort, nil)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {

	// Message lives in query string.
	// Pass it through if exists otherwise serve the page
	message := r.URL.Query().Get("chatMessageInput")

	if message == "" {
		fmt.Fprintf(w, chatPage, receivedChatMessages)
	} else {
		sendMessage(message)
		http.Redirect(w, r, "/chat", 301)
	}
}
