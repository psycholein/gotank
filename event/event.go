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

var event chan Event
var register map[string]string

func InitEvents() {
	register = make(map[string]string)
	event = make(chan Event)
	go handleEvents()
}

func SendEvent(e Event) {
	event <- e
}

func RegisterEvent(srcModule string, destModule string) {
	register[srcModule] = destModule
}

func handleEvents() {
	for e := range event {
		if d, ok := register[e.Module]; ok {
			fmt.Print("TODO", d)
		}
		fmt.Println(e)
	}
}

func Stop() {
	close(event)
}
