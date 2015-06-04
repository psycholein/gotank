package event

import "fmt"

type Event struct {
	Module  string
	Name    string
	Task    string
	Command string
}

var event chan Event

func InitEvents() {
	event = make(chan Event)
	go handleEvents()
}

func SendEvent(e Event) {
	event <- e
}

func handleEvents() {
	for e := range event {
		fmt.Println(e)
	}
}
