package main

import (
	"gotank/libs/yaml"
	"gotank/modules"
	"io/ioutil"
	"log"
)

type config struct{}

type mainConfig struct {
	Modules []string
}

func initModules() {
	m := mainConfig{}
	config{}.readConfig("./config/modules.yml", &m)

	for _, value := range m.Modules {
		modules.InitModule(value, config{})
	}
}

func (c config) readConfig(filename string, data interface{}) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(b, data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
