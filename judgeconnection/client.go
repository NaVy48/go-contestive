package judgeconnection

import (
	"fmt"
	"net"
	"time"
)

func dial(address string) (Conn, error) {
	d := net.Dialer{Timeout: 30 * time.Second}
	conn, err := d.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	tcpConn := conn.(*net.TCPConn)

	return NewConn(tcpConn), nil
}

type Client struct {
	Conn
	judgeID int
	secret  string
}

func NewClient(judgeID int, secret, address string) (Client, error) {
	conn, err := dial(address)
	if err != nil {
		return Client{}, err
	}

	conn.Write(AuthRequest{
		JudgeID: judgeID,
		Secret:  secret,
	})

	val := conn.Read()

	authAck, ok := val.(AuthResponse)
	if !ok || !authAck.OK {
		return Client{}, fmt.Errorf("auth failed")
	}

	return Client{
		conn,
		judgeID,
		secret,
	}, nil
}
