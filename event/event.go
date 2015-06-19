package event

import (
	"fmt"
)

type Event struct {
	Name   string
	Module string
	Task   string
	Value  string
}

type EventFunc func(Event)

var event chan Event
var register map[string]EventFunc

func InitEvents() {
	register = make(map[string]EventFunc)
	event = make(chan Event)
	go handleEvents()
}

func SendEvent(e Event) {
	event <- e
}

func RegisterEvent(srcModule string, destModule EventFunc) {
	register[srcModule] = destModule
}

func handleEvents() {
	for e := range event {
		if d, ok := register[e.Module]; ok {
			d(e)
		}
		fmt.Println(e)
	}
}

func Stop() {
	close(event)
}
