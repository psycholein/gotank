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
	Config() interface{}
}

type Module struct {
	Name string
	Attr moduleInterface
	Web  bool
}

var AvailableModules map[string]Module

func (m Module) Register() {
	if AvailableModules == nil {
		AvailableModules = make(map[string]Module)
	}
	AvailableModules[m.Name] = m
}

func InitModule(name string) {
	if val, ok := AvailableModules[name]; ok {
		config.Read("./config/"+val.Name+".yml", val.Attr.Config())
		val.Attr.Start()
		event.RegisterEvent(val.Name, val.Attr.GetEvent)
		return
	}
}

func StopModules() {
	for _, m := range AvailableModules {
		m.Attr.Stop()
	}
}
