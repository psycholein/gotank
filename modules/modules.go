package modules

import (
	"gotank/config"
	"gotank/event"
)

type moduleInterface interface {
	Start()
	Stop()
	GetEvent(e event.Event)
}

type Module struct {
	Name string
	Attr moduleInterface
}

var availableModules map[string]Module

func (m Module) Register(c interface{}) {
	config.Read("./config/"+m.Name+".yml", c)
	if availableModules == nil {
		availableModules = make(map[string]Module)
	}
	availableModules[m.Name] = m
}

func InitModule(name string) {
	if val, ok := availableModules[name]; ok {
		val.Attr.Start()
		return
	}
}
