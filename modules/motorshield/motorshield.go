package motorshield

import (
	"fmt"
	"gotank/components/l293d"
	"gotank/modules"
	"time"
)

const name = "motorshield"

type motorshieldModule struct{}

var running = false

func Register() {
	m := modules.Module{name, motorshieldModule{}}
	m.Register()
}

func (m motorshieldModule) Start(c interface{}) {
	go startMotor()
}

func (m motorshieldModule) Stop() {
	running = false
}

func startMotor() {
	fmt.Println("Motor Start")
	running = true
	motor := l293d.InitMotor(9, 25, 22, 23, 18, 3)
	// latch int, clk int, enable int, data int, pwm int, motor int
	motor.Forward()
	time.Sleep(1 * time.Second)
	motor.Backward()
	time.Sleep(1 * time.Second)
	motor.Stop()
	fmt.Println("Motor Stop")
	time.Sleep(10 * time.Second)
}
