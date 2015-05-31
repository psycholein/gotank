package motorshield

import (
	"fmt"
	"gotank/components/l293d"
	"gotank/modules"

	"time"
)

const name = "motorshield"

type motorshieldModule struct{}

type side struct {
	Motor, Pwm int
}
type conf struct {
	Latch, Clk, Enable, Data int
	Left, Right              side
}

var running = false
var data map[string]conf

func Register() {
	data = make(map[string]conf)
	m := modules.Module{name, motorshieldModule{}}
	m.Register(&data)
}

func (m motorshieldModule) Start() {
	fmt.Println(data)
	go startMotor()
}

func (m motorshieldModule) Stop() {
	running = false
}

func startMotor() {
	running = true
	motor := l293d.InitMotor(9, 25, 22, 23, 18, 3)
	// latch int, clk int, enable int, data int, pwm int, motor int
	for running {
		motor.Forward()
		fmt.Println("Forward")
		time.Sleep(2 * time.Second)
		motor.Backward()
		fmt.Println("Backward")
		time.Sleep(2 * time.Second)
		motor.Stop()
		fmt.Println("Stop")
		time.Sleep(2 * time.Second)
	}
}
