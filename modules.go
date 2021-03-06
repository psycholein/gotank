package main

import (
	"gotank/config"
	"gotank/event"
	"gotank/modules"
	"gotank/modules/accelerometer"
	"gotank/modules/compass"
	"gotank/modules/motion"
	"gotank/modules/motorshield"
	"gotank/modules/motorshield4wd"
	"gotank/modules/ultrasonic"
	"gotank/modules/webcam"
)

type mainConfig struct {
	Modules []string
}

func registerModules() {
	motorshield.Register()
	motorshield4wd.Register()
	ultrasonic.Register()
	motion.Register()
	compass.Register()
	accelerometer.Register()
	webcam.Register()
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
			e := event.NewEvent(module.Name, "module", "web")
			e.AddData("value", "load")
			sendData(c, e)
			for _, typ := range module.Attr.Active() {
				e := event.NewEvent(module.Name, typ, "web")
				e.AddData("value", "init")
				sendData(c, e)
			}
		}
	}
}
