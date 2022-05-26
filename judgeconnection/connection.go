package judgeconnection

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

var ErrConnClosed = fmt.Errorf("Conn closed")

type Conn interface {
	Read() Message
	Write(m Message)
	Close()
	Closed() <-chan struct{}
}

type conn struct {
	tcpConn  *net.TCPConn
	readBuf  chan Message
	writeBuf chan Message
	close    chan struct{}
	Timeout  time.Duration
}

func NewConn(tcpConn *net.TCPConn) Conn {

	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(10 * time.Second)

	log.Printf("Connected to: %s\n", tcpConn.RemoteAddr())

	c := &conn{
		tcpConn,
		make(chan Message),
		make(chan Message),
		make(chan struct{}),
		90 * time.Second,
	}

	go func(c *conn) {
		encoder := gob.NewEncoder(c.tcpConn)
		for {
			var dto DTO
			var m Message
			select {
			case m = <-c.writeBuf:
				log.Printf("Judgeconn sending: %T %v\n", m, m)
			case <-c.close:
				return
			}

			switch v := m.(type) {
			default:
				log.Printf("Unknown message type %T", m)
				return
			case AuthRequest:
				dto.AuthRequest = &v
			case AuthResponse:
				dto.AuthResponse = &v
			case SubmitRequest:
				dto.SubmitRequest = &v
			case SubmitAck:
				dto.SubmitAck = &v
			case ProblemPackageRequest:
				dto.ProblemPackageRequest = &v
			case ProblemPackageResponse:
				dto.ProblemPackageResponse = &v
			case JudgeResult:
				dto.JudgeResult = &v
			}

			err := encoder.Encode(dto)
			if err != nil {
				c.Close()
				return
			}

			log.Printf("Judgeconn sent: %T %v\n", m, m)
		}
	}(c)

	go func(c *conn) {
		dec := gob.NewDecoder(c.tcpConn)
		for {
			var dto DTO
			var m Message
			err := dec.Decode(&dto)
			if err != nil {
				c.Close()
				return
			}

			switch {
			case dto.AuthRequest != nil:
				m = *dto.AuthRequest
			case dto.AuthResponse != nil:
				m = *dto.AuthResponse
			case dto.SubmitRequest != nil:
				m = *dto.SubmitRequest
			case dto.SubmitAck != nil:
				m = *dto.SubmitAck
			case dto.ProblemPackageRequest != nil:
				m = *dto.ProblemPackageRequest
			case dto.ProblemPackageResponse != nil:
				m = *dto.ProblemPackageResponse
			case dto.JudgeResult != nil:
				m = *dto.JudgeResult
			}

			log.Printf("Judgeconn received: %T %v\n", m, m)

			select {
			case c.readBuf <- m:
			case <-c.close:
				return
			}
		}
	}(c)

	return c
}

func (c *conn) Close() {
	select {
	case <-c.close:
	default:
		close(c.close)
	}
	log.Println("Closing tcp conn")
	c.tcpConn.Close()
}

func (c *conn) Read() Message {
	c.tcpConn.SetReadDeadline(time.Now().Add(c.Timeout))
	select {
	case m := <-c.readBuf:
		c.tcpConn.SetReadDeadline(time.Time{})
		return m
	case <-c.close:
		panic(ErrConnClosed)
	}
}

func (c *conn) Write(m Message) {
	if m == nil {
		panic(fmt.Errorf("nil message is not accepted"))
	}
	select {
	case c.writeBuf <- m:
	case <-c.close:
		panic(ErrConnClosed)
	}
}

func (c *conn) Closed() <-chan struct{} {
	return c.close
}
