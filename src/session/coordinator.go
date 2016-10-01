package session

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"time"
)

var server = flag.String("server", "", "The remote server to connect to. If left blank, will listen for connections instead.")

var version = uint32(0x000100) //v00.01.00

func SetupSession() {
	if *server == "" {
		// Start listening
		ln, err := net.Listen("tcp", ":6543")
		if err != nil {
			fmt.Printf("Failed to connect: %s\r\n", err.Error())
			return
		}
		for {
			conn, err := ln.Accept()
			if err != nil {
				panic(err)
			}
			if handshake(conn, true) {
				setupSlaveSession(conn, ln)
				break
			}
		}

	} else {
		// Try connecting
		for {
			conn, err := net.Dial("tcp", *server)
			if err != nil {
				fmt.Printf("Failed to connect: %s\r\n", err.Error())
				return
			}
			if handshake(conn, true) {
				setupMasterSession(conn)
				break
			}
		}
	}
}

func setupSlaveSession(conn net.Conn, ln net.Listener) {
	pool := newPoolSession()
	session = pool

	go func() {
		for {
			req := uint8(0)
			select {
			case req = <-pool.requests:
			case <-time.After(time.Duration(60) * time.Second):
			}
			len, err := conn.Write([]byte{req})
			if len != 1 {
				panic("len != 1")
			}
			if err != nil {
				panic(err)
			}
		}
	}()

	fmt.Println("Connection successful!")

	for {
		pconn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		if !handshake(pconn, false) {
			panic("Failed to handshake pool connection")
		}
		poolType := make([]byte, 1)
		n, err := conn.Read(poolType)
		if n != 1 || err != nil {
			panic("Failed to send pool conn type")
		}
		switch poolType[0] {
		case 1:
			// Requested by remote pool
			go handleServerConnection(pool, pconn)
		case 2:
			// Requested by local pool
			pool.intake <- pconn
		default:
			panic("Impossible")
		}

	}
}

func setupMasterSession(conn net.Conn) {
	var slaveReq, quit chan uint8
	slaveReq = make(chan uint8)
	pool := newPoolSession()
	session = pool

	go func() {
		buf := make([]byte, 1)
		for {
			n, err := conn.Read(buf)
			if n != 1 || err != nil {
				quit <- 1
				break
			}
			if buf[0] > 0 {
				slaveReq <- buf[0]
			}
		}
	}()

	fmt.Println("Connection successful!")

outer:
	for {
		poolType := uint8(1)
		num := uint8(1)
		select {
		case num = <-pool.requests:
		case num = <-slaveReq:
			poolType = uint8(2)
		case <-quit:
			break outer
		}
		for i := 0; i < int(num); i++ {
			pconn, err := net.Dial("tcp", *server)
			if err != nil {
				fmt.Printf("Failed to connect: %s\r\n", err.Error())
				return
			}
			if !handshake(pconn, false) {
				panic("Failed to handshake pool connection")
			}
			n, err := conn.Write([]byte{poolType})
			if n != 1 || err != nil {
				panic("Failed to send pool conn type")
			}
			switch poolType {
			case 1:
				// Requested by local pool
				pool.intake <- pconn
			case 2:
				// Requested by remote pool
				go handleServerConnection(pool, pconn)
			default:
				panic("Impossible")
			}
		}
	}

	conn.Close()
}

func handshake(conn net.Conn, isCoord bool) bool {
	magic := []byte("Transfer")
	magicLen := 8
	isCoordByte := uint8(0)
	if isCoord {
		isCoordByte = 1
	}
	versionBlob := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBlob, version)
	blob := append(append(magic[0:magicLen], versionBlob...), isCoordByte)
	//conn.SetDeadline(time.Duration(30) * time.Second)
	n, err := conn.Write(blob)
	if n != len(blob) || err != nil {
		fmt.Println("Failed to send handshake")
		return false
	}
	n, err = conn.Read(blob)
	if n != len(blob) || err != nil {
		fmt.Println("Failed to receive handshake")
		return false
	}
	if !bytes.Equal(blob[0:magicLen], magic[0:magicLen]) {
		fmt.Println("Connection is an imposter")
		return false
	}
	if !bytes.Equal(blob[magicLen:magicLen+4], versionBlob) {
		fmt.Println("Version mismatch")
		return false
	}
	if blob[magicLen+4] != isCoordByte {
		fmt.Println("Pool/coordinator mismatch")
		return false
	}
	//conn.SetDeadline(time.Duration(75) * time.Second)
	return true
}

func handleServerConnection(session Session, conn net.Conn) {

	buf := make([]byte, 1)
	for {
		n, err := conn.Read(buf)
		if n != 1 || err != nil {
			return
		}
		if buf[0] != 0 {
			break
		}
		conn.Write([]byte{0x57})
	}
	id := int(buf[0])

	s, ok := m[id]
	if !ok {
		panic(fmt.Sprintf("No service registered for id %d", id))
	}

	s.HandleConnection(session, conn)

}
