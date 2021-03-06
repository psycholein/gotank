package motion

import (
	"gotank/event"
	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
	"gotank/modules"
)

const name = "motion"
const (
	still  = 0
	motion = 1
)

type motionModule struct{}
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
	data    = make(map[string]conf)
	status  = make(chan statusType)
	sensor  string
)

func Register() {
	m := modules.Module{name, motionModule{}, true}
	m.Register()
}

func (m motionModule) Config() interface{} {
	return &data
}

func (m motionModule) Start() {
	running = true
	go watch()
}

func (m motionModule) Stop() {
	running = false
	close(status)
}

func (m motionModule) GetEvent(e event.Event) {
}

func (m motionModule) Active() []string {
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

		err := inputs[key].Watch(embd.EdgeBoth, func(input embd.DigitalPin) {
			read, _ := input.Read()
			status <- statusType{read, input.N()}
		})
		if err != nil {
			panic(err)
		}
		e := event.NewEvent(name, key, "web")
		e.AddData("value", "start")
		e.SendEventToWeb()
	}

	for val := range status {
		state := "0"
		if val.status == motion {
			state = "1"
		}
		for key, value := range data {
			if value.Input == val.pinNum {
				e := event.NewEvent(name, key, "motion")
				e.AddData("value", state)
				e.SendEventToAll()
				break
			}
		}
	}

	for key := range data {
		inputs[key].StopWatching()
		inputs[key].Close()
		e := event.NewEvent(name, key, "web")
		e.AddData("value", "stop")
		e.SendEventToWeb()
	}
}
