package compass

import (
	"gotank/components/hmc5883l"
	"gotank/event"
	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
	"gotank/modules"
	"strconv"
	"time"
)

const name = "compass"

type compassModule struct{}
type conf struct {
	Degree, Min, Rotation int
}

var (
	running = false
	data    = make(map[string]conf)
	compass = make(map[string]*hmc5883l.HMC5883LDriver)
)

func Register() {
	m := modules.Module{name, compassModule{}, true}
	m.Register()
}

func (m compassModule) Config() interface{} {
	return &data
}

func (m compassModule) Start() {
	running = true
	go read()
}

func (m compassModule) Stop() {
	running = false
}

func (m compassModule) GetEvent(e event.Event) {
}

func (m compassModule) Active() []string {
	var active []string
	if data == nil {
		return active
	}

	for key := range data {
		active = append(active, key)
	}
	return active
}

func read() {
	for key := range data {
		compass[key] = hmc5883l.NewHMC5883LDriver(embd.NewI2CBus(1))
	}

	for running {
		for key := range data {
			degree := compass[key].Heading()
			e := event.NewEvent(name, key, "compass")
			e.AddData("value", strconv.FormatFloat(degree, 'f', 2, 64))
			e.SendEventToAll()
		}
		time.Sleep(50 * time.Millisecond)
	}
}
