package faust

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"

	flagsender "github.com/VaiTon/go-ad"
)

type Submitter struct {
	conn          net.Conn
	r             *bufio.Reader
	welcomeBanner string

	listenersMutex *sync.Mutex
	listeners      map[string]chan flagsender.Result
	listenerChan   chan any
}

func (s *Submitter) Dial(endpoint string) (*Submitter, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not connect to %s: %w", endpoint, err)
	}

	reader := bufio.NewReader(conn)

	var welcomeBanner strings.Builder
	endBanner := false
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("could not read welcome banner: %w", err)
		}

		if b == '\n' && endBanner {
			break
		} else if b == '\n' {
			endBanner = true
		} else {
			endBanner = false
		}
		welcomeBanner.WriteByte(b)
	}

	go s.incomingReceiver()

	return &Submitter{
		conn:          conn,
		r:             reader,
		welcomeBanner: welcomeBanner.String(),
	}, nil
}

func (s *Submitter) incomingReceiver() {
	for {
		select {
		case <-s.listenerChan:
			break
		default:
			line, err := s.r.ReadString('\n')
			if err != nil {
				return
			}

			line = strings.TrimSpace(line) // remove final \n

			splits := strings.Split(line, " ")
			if len(splits) < 2 {
				continue // skip invalid line
			}

			echoFlag := splits[0]
			status := splits[1]
			msg := strings.Join(splits[2:], " ")

			result := flagsender.Result{
				Success: status == "OK",
				Status:  status,
				Msg:     msg,
			}

			// is there a listener?
			s.listenersMutex.Lock()
			listener, ok := s.listeners[echoFlag]
			s.listenersMutex.Unlock()

			if !ok {
				slog.Error("illegal condition")
				continue
			}

			listener <- result
		}

	}
}

func (s *Submitter) Send(flag string) (flagsender.Result, error) {

	// add listener
	resChan := make(chan flagsender.Result)
	s.listenersMutex.Lock()
	s.listeners[flag] = resChan
	s.listenersMutex.Unlock()

	// write string
	_, err := s.conn.Write([]byte(flag + "\n"))
	if err != nil {
		return flagsender.Result{}, fmt.Errorf("could not send flag: %w", err)
	}

	// wait for listener
	return <-resChan, nil
}

func (s *Submitter) Close() error {
	return s.conn.Close()
}
