package motorshield

import (
	"gotank/components/l293d"
	"gotank/event"
	"gotank/modules"
)

const name = "motorshield"

type motorshieldModule struct{}

type side struct {
	Motor, Pwm int
}
type conf struct {
	Latch, Clk, Enable, Data int
	Left, Right              side
}

var (
	running = false
	data    = make(map[string]conf)
	left    = make(map[string]*l293d.MotorL293d)
	right   = make(map[string]*l293d.MotorL293d)
)

func Register() {
	m := modules.Module{name, motorshieldModule{}, true}
	m.Register()
}

func (m motorshieldModule) Config() interface{} {
	return &data
}

func (m motorshieldModule) Start() {
	go startMotor()

	event.RegisterEvent("ultrasonic", m.GetEvent) // Test
}

func (m motorshieldModule) Stop() {
	running = false
	for key := range data {
		e := event.NewEvent(name, key, "web")
		e.AddData("value", "stop")
		e.SendEventToWeb()
	}
}

func (m motorshieldModule) GetEvent(e event.Event) {
	if e.Module != name {
		return
	}
	if e.Task == "control" {
		c := control{e.Name, e.Data}
		c.handleControl()
	}
	if e.Task == "speed" {
		handleSpeed(e.Name, e.Data)
	}
}

func (m motorshieldModule) Active() []string {
	var active []string
	if data == nil {
		return active
	}

	for key := range data {
		active = append(active, key)
	}
	return active
}

func startMotor() {
	running = true
	for key, value := range data {
		// latch int, clk int, enable int, data int - pwm int, motor int
		l := l293d.InitL293d(value.Latch, value.Clk, value.Enable, value.Data)
		left[key] = l.InitMotor(value.Left.Pwm, value.Left.Motor)
		right[key] = l.InitMotor(value.Right.Pwm, value.Right.Motor)
		e := event.NewEvent(name, key, "web")
		e.AddData("value", "start")
		e.SendEventToWeb()
	}
}
