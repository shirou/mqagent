package main

import (
	"encoding/json"
)

type ActionJson struct {
	ActionId string   `json:"actionid"`
	Action   string   `json:"action"`
	Args     []string `json:"args"`
	Type     string   `json:"type"`
	Interval int      `json:"interval"`
	Command  string   `json:"command"`
	Handlers []string `json:"handlers"`

	TopicRoot string
}

func (a ActionJson) String() string {
	s, _ := json.Marshal(a)
	return string(s)
}
