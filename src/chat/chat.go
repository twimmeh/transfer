package chat

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"session"
	"strconv"
	"strings"
)

const chatServiceId = 1

var receivedChatMessages = ""
var chatPort = flag.String("chatport", "8080", "The HTTP port of the chat service. Default: 8080")

const chatPage = `<!DOCTYPE html>
<html>
<head>
<title>Chatz</title>
  <style>
    textarea {
        resize: none;
    }
  </style>
</head>
<body>
<!-- Start your code here -->
  <div>
   <textarea id="inputText" rows="4" cols="50" autofocus="true"></textarea>
  </div>
  <div>
<button onclick="writeMessage();" type="button" id="sendButton">Send</button>
  </div>
  <div>
  <textarea rows="4" cols="50" readonly="true" id="outputText"></textarea>
  </div>
  
  <script>

function writeMessage(){
	var inputText = document.getElementById("inputText");

	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/chat/send", true);
	xhr.setRequestHeader("Content-Type", "text/plain; charset=UTF-8");
	xhr.send(inputText.value);

	inputText.value = "";
	};
  </script>
</body>
</html>
`

//Register chat service
func init() {
	chatService := ChatService{}
	session.RegisterService(chatServiceId, chatService)
}

type ChatService struct{}

func (svc ChatService) HandleConnection(session session.Session, conn session.Connection) {

	// Dump to console (for now)
	b, _ := ioutil.ReadAll(conn)
	fmt.Println(string(b[:]))
}

// If a session is available, try to send a message to it
func sendMessage(message string) {
	s, _ := session.GetSession()
	if s != nil {
		conn := s.OpenConnection(chatServiceId)
		io.WriteString(conn, message)
		conn.Close()
		fmt.Println("Message sent!")
	} else {
		fmt.Println("No chat sessions available")
	}
}

// Web-based chat message sending
func SetupWebService() string {

	if _, err := strconv.Atoi(*chatPort); err != nil {
		panic("Invalid chat port specified (must be 32-bit integer). You put: " + *chatPort)
	}

	http.HandleFunc("/chat", chatHandler)      // Main chat page
	http.HandleFunc("/chat/send", sendHandler) // Endpoint for sending messages
	http.HandleFunc("/chat/read", readHandler) // Endpoint for receiving messages (not yet implemented)
	return *chatPort
}

// Serve the chat page
func chatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, chatPage)
}

// Parse request for chat message and send to connected client
func sendHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodPost:
		contentType := r.Header.Get("Content-Type")

		if strings.HasPrefix(contentType, "text/plain") {

			//todo: handle edge cases more gracefully, e.g. no remote client connected.
			b, _ := ioutil.ReadAll(r.Body)
			fmt.Println("Trying to send message: " + string(b[:]))
			sendMessage(string(b[:]))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Content-Type must be text/plain")
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "You must POST chat messages.")
	}
}

func readHandler(w http.ResponseWriter, r *http.Request) {

	// todo!

}
