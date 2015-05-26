package modules

import (
	"fmt"
	"gotank/modules/motorshield"
)

type motorshieldModule struct {
	Name string
}

func init() {
	name := motorshield.Name
	m := module{name, motorshieldModule{name}}
	m.register()
}

func (m motorshieldModule) Start(c interface{}) {
	fmt.Println(m.Name, c)
}

func (m motorshieldModule) Stop() {
}
