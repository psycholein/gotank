package ultrasonic

import "gotank/modules"

const name = "ultrasonic"

type ultrasonicModule struct {
	Name string
}

func Register() {
	m := modules.Module{name, ultrasonicModule{name}}
	m.Register()
}

func (m ultrasonicModule) Start(c interface{}) {
}

func (m ultrasonicModule) Stop() {
}
