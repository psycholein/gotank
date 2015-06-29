package motionsensor

import (
	"gotank/event"
	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
	"gotank/modules"
)

const name = "motionsensor"
const (
	still  = 0
	motion = 1
)

type motionsensorModule struct{}
type position struct {
	Distance, Degree int
}
type conf struct {
	Input    int
	Position position
}
type statusType struct {
	status int
	pinNum int
}

var (
	running = false
	data    map[string]conf
	status  chan statusType
	sensor  string
)

func Register() {
	data = make(map[string]conf)
	m := modules.Module{name, motionsensorModule{}, true}
	m.Register()
}

func (m motionsensorModule) Config() interface{} {
	return &data
}

func (m motionsensorModule) Start() {
	running = true
	go watch()
}

func (m motionsensorModule) Stop() {
	running = false
	close(status)
}

func (m motionsensorModule) GetEvent(e event.Event) {
}

func (m motionsensorModule) Active() []string {
	var active []string
	if data == nil {
		return active
	}

	for key := range data {
		active = append(active, key)
	}
	return active
}

func watch() {
	inputs := make(map[string]embd.DigitalPin)
	for key, value := range data {
		inputs[key], _ = embd.NewDigitalPin(value.Input)
		inputs[key].SetDirection(embd.In)

		status = make(chan statusType, 1)
		err := inputs[key].Watch(embd.EdgeBoth, func(input embd.DigitalPin) {
			input.N()
			read, _ := input.Read()
			status <- statusType{read, input.N()}
		})
		if err != nil {
			panic(err)
		}
		event.SendEvent(event.Event{name, key, "web", "start"})
	}

	for val := range status {
		state := "0"
		if val.status == motion {
			state = "1"
		}
		for key, value := range data {
			if value.Input == val.pinNum {
				event.SendEvent(event.Event{name, key, "motion", state})
				break
			}
		}
	}

	for key := range data {
		inputs[key].StopWatching()
		inputs[key].Close()
		event.SendEvent(event.Event{name, key, "web", "stop"})
	}
}
