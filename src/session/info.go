package session

import (
	"flag"
	"net"
)

var name = flag.String("name", "", "The name to identify this transfer session.")

// SessionInfo provides info about the session.
type SessionInfo struct {
	LocalName  string
	RemoteName string
}

func exchangeSessionInfo(conn net.Conn) SessionInfo {
	var nameToSend, tmpLocalName, tmpRemoteName string
	tmpLocalName = *name
	if len(tmpLocalName) > 255 {
		tmpLocalName = ""
	}
	if tmpLocalName == "" {
		tmpLocalName, nameToSend = "You", "Other"
	} else {
		tmpLocalName, nameToSend = tmpLocalName, tmpLocalName
	}
	length, err := conn.Write([]uint8{uint8(len(nameToSend))})
	if err != nil || length != 1 {
		panic(err)
	}
	length, err = conn.Write([]byte(nameToSend))
	if err != nil || length != len(nameToSend) {
		panic(err)
	}
	buf := make([]byte, 1)
	length, err = conn.Read(buf)
	if err != nil || length != 1 {
		panic(err)
	}
	rlength := int(buf[0])
	buf = make([]byte, rlength)
	length, err = conn.Read(buf)
	if err != nil || length != rlength {
		panic(err)
	}
	tmpRemoteName = string(buf)
	return SessionInfo{tmpLocalName, tmpRemoteName}
}
