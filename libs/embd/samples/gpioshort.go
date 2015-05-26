// +build ignore

package main

import (
	"flag"

	"gotank/libs/embd"

	_ "gotank/libs/embd/host/all"
)

func main() {
	flag.Parse()

	embd.InitGPIO()
	defer embd.CloseGPIO()

	embd.SetDirection(10, embd.Out)
	embd.DigitalWrite(10, embd.High)
}
