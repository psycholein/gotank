package modules

import (
	"fmt"
	//_ "gotank/libs/embd"
	//_ "gotank/libs/embd/host/all" // select the right board
)

type module struct {
	Name string
}

type event struct {
	Name string
	Data string
}

type moduleInterface interface {
	Start()
	Stop()
	Events() []string
	Event(e event)
}

var modules map[string]module

func (m module) register() {
	if modules == nil {
		modules = make(map[string]module)
	}
	modules[m.Name] = m
}

func InitModule(name string) {
	if val, ok := modules[name]; ok {
		fmt.Println(val)
		return
	}
}
