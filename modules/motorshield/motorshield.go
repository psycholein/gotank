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
var left, right map[string]l293d.MotorShieldL293d

func Register() {
	data = make(map[string]conf)
	m := modules.Module{name, motorshieldModule{}}
	m.Register(&data)
}

func (m motorshieldModule) Start() {
	fmt.Println(data)
	left := make(map[string]l293d.MotorShieldL293d)
	right := make(map[string]l293d.MotorShieldL293d)
	for key, value := range data {
		// latch int, clk int, enable int, data int, pwm int, motor int
		left[key] = *l293d.InitMotor(value.Latch, value.Clk, value.Enable, value.Data, value.Left.Pwm, value.Left.Motor)
		right[key] = *l293d.InitMotor(value.Latch, value.Clk, value.Enable, value.Data, value.Right.Pwm, value.Right.Motor)
	}
	go startMotor()
}

func (m motorshieldModule) Stop() {
	running = false
}

func startMotor() {
	running = true
	for running {
		for key := range data {
			left[key].Forward()
			fmt.Println("Forward left")
			time.Sleep(1 * time.Second)
			left[key].Backward()
			fmt.Println("Backward left")
			time.Sleep(1 * time.Second)
			left[key].Stop()

			right[key].Forward()
			fmt.Println("Forward right")
			time.Sleep(1 * time.Second)
			right[key].Backward()
			fmt.Println("Backward right")
			time.Sleep(1 * time.Second)
			right[key].Stop()
			fmt.Println("Stop")
			time.Sleep(2 * time.Second)
		}
	}
}
