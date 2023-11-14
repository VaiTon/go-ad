package saarctf

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatus(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, `{
		"teams": [
			{
				"id": 1,
				"name": "NOP",
				"ip": "10.32.1.2"
			},
			{
				"id": 2,
				"name": "saarsec",
				"ip": "10.32.2.2"
			}
		],
		"flag_ids": {
			"service_1": {
				"10.32.1.2": {
					"15": ["username1", "username1.2"],
					"16": ["username2", "username2.2"]
				},
				"10.32.2.2": {
					"15": ["username3", "username3.2"],
					"16": ["username4", "username4.2"]
				}
			}
		}
	}`)
	}))

	res, err := GetStatus(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}
