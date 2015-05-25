package main

import (
	"gotank/libs/yaml"
	"gotank/modules"
	"io/ioutil"
	"log"
)

type mainConfig struct {
	Modules []string
}

func initModules() {
	m := mainConfig{}
	b, err := ioutil.ReadFile("./config/modules.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(b, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, value := range m.Modules {
		modules.InitModule(value)
	}
}
