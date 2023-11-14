package saarctf

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/VaiTon/go-ad"
)

type Sender struct {
	conn *net.TCPConn
	rw   *bufio.ReadWriter
}

func NewSender(endpoint *net.TCPAddr) (*Sender, error) {
	conn, err := net.DialTCP("tcp", nil, endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not connect to %s: %w", endpoint, err)
	}

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	return &Sender{conn: conn, rw: rw}, nil
}

func (s *Sender) Send(flag string) (flagsender.Result, error) {
	flag += "\n" // "Each flag must be submitted in a single line terminated by a line feed (\n)."

	_, err := s.rw.WriteString(flag)
	if err != nil {
		return flagsender.Result{}, fmt.Errorf("could not send flag: %w", err)
	}

	resp, err := s.rw.ReadString('\n')
	if err != nil {
		return flagsender.Result{}, fmt.Errorf("could not read response: %w", err)
	}

	splits := strings.Split(resp, " ")
	if len(splits) < 1 {
		return flagsender.Result{}, fmt.Errorf("invalid response: %s", resp)
	}

	status := splits[0]
	status = strings.Trim(status, "[]")

	result := flagsender.Result{Status: status, Msg: resp, Success: false}
	if status == "OK" {
		result.Success = true
	}

	return result, nil
}

func (s *Sender) Close() error {
	return s.conn.Close()
}
