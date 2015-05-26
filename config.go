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
	c := config{}
	c.readConfig("./config/modules.yml", &m)

	for _, value := range m.Modules {
		i := make(map[interface{}]interface{})
		c.readConfig("./config/"+value+".yml", &i)
		modules.InitModule(value, i)
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
