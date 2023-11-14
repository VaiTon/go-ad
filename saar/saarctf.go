package saar

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/VaiTon/go-ad"
)

type Submitter struct {
	conn net.Conn
	r    *bufio.Reader
}

func Dial(endpoint string) (*Submitter, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not connect to %s: %w", endpoint, err)
	}

	return &Submitter{
		conn: conn,
		r:    bufio.NewReader(conn),
	}, nil
}

func (s *Submitter) Send(flag string) (flagsender.Result, error) {
	_, err := s.conn.Write([]byte(flag + "\n"))
	if err != nil {
		return flagsender.Result{}, fmt.Errorf("could not send flag: %w", err)
	}

	resp, err := s.r.ReadString('\n')
	if err != nil {
		return flagsender.Result{}, fmt.Errorf("could not read response: %w", err)
	}

	resp = strings.TrimSpace(resp) // remove final \n

	splits := strings.Split(resp, " ")
	if len(splits) < 1 {
		return flagsender.Result{}, fmt.Errorf("invalid response: %s", resp)
	}

	status := splits[0]
	status = strings.TrimPrefix(status, "[")
	status = strings.TrimSuffix(status, "]")

	return flagsender.Result{
		Status:  status,
		Msg:     resp,
		Success: status == "OK",
	}, nil
}

func (s *Submitter) Close() error {
	return s.conn.Close()
}
