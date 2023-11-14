package saarctf

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	Service       map[Round][]string
	Round         int
	Team          string
	ServiceFlagId map[Team]Service
	Status        struct {
		Teams   []TeamInfo               `json:"teams"`
		FlagIds map[string]ServiceFlagId `json:"flag_ids"`
	}

	TeamInfo struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Ip   string `json:"ip"`
	}
)

func GetStatus(statusUrl string) (Status, error) {
	res, err := http.Get(statusUrl)
	if err != nil {
		return Status{}, fmt.Errorf("could not get status: %w", err)
	}

	var status Status
	err = json.NewDecoder(res.Body).Decode(&status)
	if err != nil {
		return Status{}, fmt.Errorf("coult not decode status: %w", err)
	}

	return status, nil
}
