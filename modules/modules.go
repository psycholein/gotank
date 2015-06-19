package modules

import (
	"gotank/config"
	"gotank/event"
)

type moduleInterface interface {
	Start()
	Stop()
	GetEvent(e event.Event)
	Active() []string
}

type Module struct {
	Name string
	Attr moduleInterface
	Web  bool
}

var AvailableModules map[string]Module

func (m Module) Register(c interface{}) {
	config.Read("./config/"+m.Name+".yml", c)
	if AvailableModules == nil {
		AvailableModules = make(map[string]Module)
	}
	AvailableModules[m.Name] = m
	if m.Web {
		event.SendEvent(event.Event{m.Name, "module", "web", "register"})
	}
}

func InitModule(name string) {
	if val, ok := AvailableModules[name]; ok {
		val.Attr.Start()
		return
	}
}

func StopModules() {
	for _, m := range AvailableModules {
		m.Attr.Stop()
	}
}
