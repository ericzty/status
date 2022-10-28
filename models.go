package main

import (
	"sync"
)

type target struct {
	Name      string
	Target    string
	Method    string
	Match     string
	Public    bool
	Status4   targetStatus
	Status6   targetStatus
	LastCheck string
}

type targetStatus struct {
	Up       bool
	Fails    uint8
	Alerting bool
}

type status struct {
	Mutex  *sync.Mutex
	values []*target
}

type embed struct {
	Content interface{} `json:"content"`
	Embeds  []embedBody `json:"embeds"`
}

type embedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

type embedBody struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Color       int         `json:"color"`
	Footer      embedFooter `json:"footer"`
}

type pushNotif struct {
	Token    string `json:"token"`
	User     string `json:"user"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
	Url      string `json:"url"`
}
