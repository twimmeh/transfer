package chat

import (
	"session"
	"fmt"
	"io/ioutil"
	"io"
)

const chatServiceId = 1

//Register chat service
func init(){
	chatService := ChatService{}
	session.RegisterService(chatServiceId, chatService)
}

type ChatService struct{}

// Simple read and dump to command line
func (svc ChatService)HandleConnection(session session.Session, conn session.Connection){

	b,_ := ioutil.ReadAll(conn)
	fmt.Println(string(b[:]))
}

// If a session is available, try to send a message to it
func SendMessage(message string){
	
	s,_ := session.GetSession()
	if(s != nil){
		conn := s.OpenConnection(chatServiceId)
		io.WriteString(conn,message)
		conn.Close()
	}else{
		fmt.Println("No chat sessions available")
	}
}