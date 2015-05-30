package ultrasonic

import "gotank/modules"

const name = "ultrasonic"

type ultrasonicModule struct{}

func Register() {
	m := modules.Module{name, ultrasonicModule{}}
	m.Register()
}

func (m ultrasonicModule) Start(c interface{}) {
}

func (m ultrasonicModule) Stop() {
}
