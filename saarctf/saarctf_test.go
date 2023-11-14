package saarctf

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

func TestNewSender(t *testing.T) {

	const addr = "localhost:1234"
	const flag = "testflag12345"

	fakeSocket, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := fakeSocket.Accept()
		if err != nil {
			t.Error(err)
			return
		}

		rd := bufio.NewReader(conn)
		recvFlag, err := rd.ReadString('\n')
		recvFlag = strings.TrimSpace(recvFlag)

		if err != nil {
			t.Error(err)
			return
		}
		if string(recvFlag) != flag {
			t.Errorf("wrong flag: %s", recvFlag)
			return
		}

		_, err = conn.Write([]byte("[OK]\n"))
		if err != nil {
			t.Error(err)
			return
		}
	}()

	sender, err := Dial(addr)
	if err != nil {
		t.Fatal(err)
	}

	res, err := sender.Send(flag)
	if err != nil {
		t.Fatal(err)
	}

	if !res.Success {
		t.Error("result is not success")
	}
	if res.Status != "OK" {
		t.Error("status is not ok:", res.Status)
	}
}
