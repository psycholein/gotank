package modules

import "gotank/modules/ultrasonic"

type ultrasonicModule struct {
	Name string
}

func init() {
	name := ultrasonic.Name
	m := module{name, ultrasonicModule{name}}
	m.register()
}

func (m ultrasonicModule) Start(c interface{}) {
	ultrasonic.Start(c)
}

func (m ultrasonicModule) Stop() {
	ultrasonic.Stop()
}
