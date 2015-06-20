package main

import (
	"gotank/config"
	"gotank/event"
	"gotank/modules"
	"gotank/modules/motorshield"
	"gotank/modules/ultrasonic"
)

type mainConfig struct {
	Modules []string
}

func registerModules() {
	motorshield.Register()
	ultrasonic.Register()
}

func initModules() {
	m := mainConfig{}
	config.Read("./config/modules.yml", &m)
	for _, value := range m.Modules {
		modules.InitModule(value)
	}
}

func stopModules() {
	modules.StopModules()
}

func sendModulesToWeb(c *connection) {
	for _, module := range modules.AvailableModules {
		if module.Web {
			sendData(c, event.Event{module.Name, "module", "web", "load"})
			for _, typ := range module.Attr.Active() {
				sendData(c, event.Event{module.Name, typ, "web", "init"})
			}
		}
	}
}
