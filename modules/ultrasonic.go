package modules

import (
	"fmt"
	"gotank/modules/ultrasonic"
)

type ultrasonicModule struct {
	Name string
}

func init() {
	name := ultrasonic.Name
	m := module{name, ultrasonicModule{name}}
	m.register()
}

func (m ultrasonicModule) Start(c interface{}) {
	fmt.Println(m.Name, c)
}

func (m ultrasonicModule) Stop() {
}
