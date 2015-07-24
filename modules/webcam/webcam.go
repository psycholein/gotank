package webcam

import (
	"fmt"
	"gotank/event"
	"gotank/modules"
)

const name = "webcam"

type webcamModule struct{}
type conf struct {
	Device string
}

var (
	running = false
	data    = make(map[string]conf)
)

func Register() {
	m := modules.Module{name, webcamModule{}, true}
	m.Register()
}

func (m webcamModule) Config() interface{} {
	fmt.Println("Read webcam config")
	return &data
}

func (m webcamModule) Start() {
	running = true
	for key := range data {
		e := event.NewEvent(name, key, "web")
		e.AddData("start", "1")
		e.SendEventToWeb()
	}
}

func (m webcamModule) Stop() {
	running = false
	for key := range data {
		e := event.NewEvent(name, key, "web")
		e.AddData("stop", "1")
		e.SendEventToWeb()
	}
}

func (m webcamModule) GetEvent(e event.Event) {
	if e.Task == "get" {
		sendImg(e.Name)
	}
}

func (m webcamModule) Active() []string {
	var active []string
	if data == nil {
		return active
	}

	for key := range data {
		active = append(active, key)
	}
	return active
}

func sendImg(key string) {
	// TODO
}
