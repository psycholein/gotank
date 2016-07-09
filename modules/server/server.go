package server

import (
	"gotank/event"
	"gotank/modules"
	"time"
)

const (
	name  = "server"
	delay = 5
)

type serverModule struct{}
type conf struct {
	Addr string
}

var (
	running = false
	data    = make(map[string]conf)
)

func Register() {
	m := modules.Module{name, serverModule{}, false}
	m.Register()
}

func (m serverModule) Config() interface{} {
	return &data
}

func (m serverModule) Start() {
	running = true
	time.Sleep(delay * time.Second)
	for key := range data {
		e := event.NewEvent(name, key, "web")
		e.AddData("start", "1")
		e.SendEventToWeb()
	}
}

func (m serverModule) Stop() {
	running = false
	for key := range data {
		e := event.NewEvent(name, key, "web")
		e.AddData("stop", "1")
		e.SendEventToWeb()
	}
}

func (m serverModule) GetEvent(e event.Event) {
}

func (m serverModule) Active() []string {
	var active []string
	if data == nil {
		return active
	}

	for key := range data {
		active = append(active, key)
	}
	return active
}
