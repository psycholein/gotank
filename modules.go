package main

import (
	"gotank/modules"
	"gotank/modules/motorshield"
	"gotank/modules/ultrasonic"
)

func registerModules() {
	motorshield.Register()
	ultrasonic.Register()
}

func initModules() {
	m := mainConfig{}
	readConfig("./config/modules.yml", &m)
	for _, value := range m.Modules {
		i := make(map[interface{}]interface{})
		readConfig("./config/"+value+".yml", &i)
		modules.InitModule(value, i)
	}
}
