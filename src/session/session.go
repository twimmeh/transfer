package session

import (
	"errors"
	"io"
	"net"
)

var m = make(map[int]Service)

var session Session

func init() {
	session = nil
}

// Session represents an active pairing of two Transfer instances, and allows
// clients (in one instance) to open a connection to services in the other.
type Session interface {
	OpenConnection(id int) Connection
}

// Connection wraps a TCP connection between a client and service.
type Connection interface {
	io.ReadWriteCloser
}

// Service is the interface that should be implemented by all services in the
// app.
//
// Every service should create an implementation of Service, and register it
// using RegisterService.
type Service interface {
	HandleConnection(session Session, connection Connection)
}

// GetSession returns the current session, or an error if one is not available.
// For initial testing purposes, GetSession just returns a loopback service.
func GetSession() (Session, error) {
	if session == nil {
		return nil, errors.New("Session not available yet")
	}
	return session, nil
}

//
// Every service should create an implementation of Service, and register it
// using RegisterService.
func RegisterService(id int, service Service) {
	if id < 0 || id > 255 {
		panic("Invalid ID")
	}
	m[id] = service
}

// Loopback service for initial testing.
// TODO: Delete this.
type loopbackService struct{}

func (l *loopbackService) OpenConnection(id int) Connection {
	if id < 0 || id > 255 {
		panic("Invalid ID")
	}
	s, ok := m[id]
	if !ok {
		panic("Service ID not found")
	}
	conn1, conn2 := net.Pipe()
	go s.HandleConnection(l, conn1)
	return conn2
}

// For initial go environment testing.
// TODO: Delete this.
func GetTestString() string {
	return "It works!"
}
