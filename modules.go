package main

import (
	"gotank/modules/motorshield"
	"gotank/modules/ultrasonic"
)

func registerModules() {
	motorshield.Register()
	ultrasonic.Register()
}
