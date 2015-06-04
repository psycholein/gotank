package main

import (
	"gotank/config"
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
