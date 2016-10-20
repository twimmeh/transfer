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
<button onclick="writeToWindow();" type="button" id="sendButton">Send</button>
  </div>
  <div>
  <textarea rows="4" cols="50" readonly="true" id="outputText"></textarea>
  </div>
  
  <script>
function writeToWindow(){
      var inputText = document.getElementById("inputText");
      var outputText = document.getElementById("outputText");
      var msg = inputText.value;
      inputText.value = null;
      
      if(outputText.value === "") {
        outputText.value = msg;
      } else {
        outputText.value += "\n"+msg;
      }
      //update cursor on outputText
      outputText.focus();
      //put focus back on input box
      inputText.focus();
    }
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

// Simple read and dump to command line
func (svc ChatService) HandleConnection(session session.Session, conn session.Connection) {
	b, _ := ioutil.ReadAll(conn)
	receivedChatMessages = receivedChatMessages + "<br>" + string(b[:])
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
func SetupWebService() (string) {

	if _, err := strconv.Atoi(*chatPort); err != nil {
		panic("Invalid chat port specified (must be 32-bit integer). You put: " + *chatPort)
	}

	http.HandleFunc("/chat", chatHandler)
	return *chatPort
}

func chatHandler(w http.ResponseWriter, r *http.Request) {

	// Message lives in query string.
	// Pass it through if exists otherwise serve the page
	message := r.URL.Query().Get("chatMessageInput")

	if message == "" {
		//fmt.Fprintf(w, chatPage, receivedChatMessages)
		fmt.Fprintf(w, chatPage)
	} else {
		sendMessage(message)
		http.Redirect(w, r, "/chat", 301)
	}
}
