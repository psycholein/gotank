package motorshield

import (
	_ "gotank/components/l293d"
	"gotank/modules"
)

const name = "motorshield"

type motorshieldModule struct{}

func Register() {
	m := modules.Module{name, motorshieldModule{}}
	m.Register()
}

func (m motorshieldModule) Start(c interface{}) {
}

func (m motorshieldModule) Stop() {
}
