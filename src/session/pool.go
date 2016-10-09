package session

import (
	"net"
	"time"
)

type poolSession struct {
	conns       chan clientConnection
	pool        chan clientConnection
	requests    chan uint8
	intake      chan net.Conn
	pending     chan uint8
	sessionInfo SessionInfo
}

func newPoolSession(conn net.Conn) *poolSession {
	p := &poolSession{
		conns:       make(chan clientConnection),
		pool:        make(chan clientConnection, 3),
		requests:    make(chan uint8, 1),
		intake:      make(chan net.Conn, 1),
		pending:     make(chan uint8, 3),
		sessionInfo: exchangeSessionInfo(conn)}

	go func() {
		for {
			shortfall := 3 - (len(p.pool) + len(p.pending))
			if shortfall > 0 {
				for i := 0; i < shortfall; i++ {
					p.pending <- 1
				}
				p.requests <- uint8(shortfall)
			}
			p.conns <- <-p.pool
		}
	}()

	go func() {
		for {
			newIntake := <-p.intake
			p.pool <- createClientConnection(newIntake)
			<-p.pending
		}
	}()

	return p
}

func (s *poolSession) OpenConnection(id int) Connection {
	if id <= 0 || id > 255 {
		panic("Invalid ID")
	}
	cc := <-s.conns
	return cc.connectService(uint8(id))
}

type clientConnection struct {
	id           chan uint8
	preparedConn chan net.Conn
}

func createClientConnection(conn net.Conn) clientConnection {
	ret := clientConnection{make(chan uint8, 0), make(chan net.Conn)}
	go func() {
		pings := make(chan int, 1)
		for {
			select {
			case id := <-ret.id:
				pings <- 1
				conn.Write([]byte{id})
				ret.preparedConn <- conn
				<-pings
				return
			case <-time.After(60):
				pings <- 1
				conn.Write([]byte{0})
				go func() {
					buf := []byte{0}
					n, err := conn.Read(buf)
					if err != nil {
						panic(err)
					}
					if n != 1 {
						panic("Pong recv failed")
					}
					if buf[0] != 0x57 {
						panic("Invalid pong")
					}
					<-pings
				}()
			}
		}
	}()
	return ret
}

func (c *clientConnection) connectService(id uint8) net.Conn {
	c.id <- id
	return <-c.preparedConn
}

func (s *poolSession) GetInfo() SessionInfo {
	return s.sessionInfo
}
