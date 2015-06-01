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
	left := l293d.InitMotor(9, 25, 22, 23, 18, 3)
	right := l293d.InitMotor(9, 25, 22, 23, 18, 4)
	// latch int, clk int, enable int, data int, pwm int, motor int
	for running {
		left.Forward()
		fmt.Println("Forward left")
		time.Sleep(1 * time.Second)
		left.Backward()
		fmt.Println("Backward left")
		time.Sleep(1 * time.Second)
		left.Stop()

		right.Forward()
		fmt.Println("Forward right")
		time.Sleep(1 * time.Second)
		right.Backward()
		fmt.Println("Backward right")
		time.Sleep(1 * time.Second)
		right.Stop()
		fmt.Println("Stop")
		time.Sleep(2 * time.Second)
	}
}
