package event

import (
	"fmt"
)

type Event struct {
	Module string
	Name   string
	Task   string
	Value  string
}

type EventFunc func(Event)

var event chan Event
var register map[string][]EventFunc

func InitEvents() {
	register = make(map[string][]EventFunc)
	event = make(chan Event)
	go handleEvents()
}

func SendEvent(e Event) {
	event <- e
}

func RegisterEvent(srcModule string, destModule EventFunc) {
	register[srcModule] = append(register[srcModule], destModule)
}

func handleEvents() {
	for e := range event {
		if items, ok := register[e.Module]; ok {
			for _, item := range items {
				item(e)
			}
		}
		if items, ok := register["_all"]; ok {
			for _, item := range items {
				item(e)
			}
		}
		fmt.Println("handleEvents", e)
	}
}

func Stop() {
	close(event)
}
