package main

import (
	"gotank/libs/yaml"
	"io/ioutil"
	"log"
)

type mainConfig struct {
	Modules []string
}

func readConfig(filename string, data interface{}) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(b, data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
