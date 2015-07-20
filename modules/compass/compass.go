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

type magnetoModule struct{}
type conf struct {
	Degree, Min, Rotation int
}

var (
	running = false
	data    map[string]conf
	compass map[string]*hmc5883l.HMC5883LDriver
)

func Register() {
	data = make(map[string]conf)
	m := modules.Module{name, magnetoModule{}, true}
	m.Register()
}

func (m magnetoModule) Config() interface{} {
	return &data
}

func (m magnetoModule) Start() {
	running = true
	go read()
}

func (m magnetoModule) Stop() {
	running = false
}

func (m magnetoModule) GetEvent(e event.Event) {
}

func (m magnetoModule) Active() []string {
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
	compass = make(map[string]*hmc5883l.HMC5883LDriver)
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
