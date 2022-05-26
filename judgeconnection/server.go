package judgeconnection

import (
	"log"
	"net"
)

type JudgeManager interface {
	RunJudge(conn Conn) error
}

type Server struct {
	jm          JudgeManager
	portAddress string
	ln          net.Listener
}

func NewJudgeServer(jm JudgeManager, portAddress string) Server {
	return Server{jm, portAddress, nil}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		if pan := recover(); pan != nil {
			log.Printf("Panic encoutered: %v", pan)
		}
	}()
	tcp, ok := conn.(*net.TCPConn)
	if !ok {
		return
	}
	defer tcp.Close()

	jconn := NewConn(tcp)
	err := s.jm.RunJudge(jconn)
	log.Printf("Judge failed and disconecting, encountered err: %v", err)
}

func (s *Server) Listen() error {
	log.Printf("Listening for judges on: %s\n", s.portAddress)
	ln, err := net.Listen("tcp", s.portAddress)
	if err != nil {
		return err
	}
	s.ln = ln
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			log.Printf("Connection acception error: %v", err)
			return err
		}
		go s.handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
	}
}

func (s *Server) Close() {
	if s.ln != nil {
		err := s.ln.Close()
		if err != nil {
			panic(err)
		}
		s.ln = nil
	}
}
