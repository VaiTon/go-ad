package saarctf

import (
	"fmt"
	"testing"
)

func TestGetStatus(t *testing.T) {
	res, err := GetStatus()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}
