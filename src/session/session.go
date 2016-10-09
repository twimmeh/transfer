package session

import (
	"errors"
	"io"
)

var m = make(map[int]Service)

var session Session

func init() {
	session = nil
}

// Session represents an active pairing of two Transfer instances.
type Session interface {
	// Make a connection to a service to the other Transfer instance.
	OpenConnection(id int) Connection

	// Returns a copy of the SessionInfo object describing useful info about this session.
	GetInfo() SessionInfo
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
func GetSession() (Session, error) {
	if session == nil {
		return nil, errors.New("Session not available yet")
	}
	return session, nil
}

// Every service should create an implementation of Service, and register it
// using RegisterService.
func RegisterService(id int, service Service) {
	if id < 0 || id > 255 {
		panic("Invalid ID")
	}
	m[id] = service
}
