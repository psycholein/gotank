package ultrasonic

import (
	"fmt"
	"gotank/modules"
	"time"

	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
)

const name = "ultrasonic"

type ultrasonicModule struct{}

var running = false

func Register() {
	m := modules.Module{name, ultrasonicModule{}}
	m.Register()
}

func (m ultrasonicModule) Start(c interface{}) {
	go start()
}

func (m ultrasonicModule) Stop() {
	running = false
}

func start() {
	running = true
	for running {
		fmt.Println(distance(17, 4))
		time.Sleep(100 * time.Millisecond)
	}
}

func distance(echoNum int, triggerNum int) (float64, error) {
	fmt.Println("Messe")
	echo, _ := embd.NewDigitalPin(echoNum)
	echo.SetDirection(embd.In)
	trigger, _ := embd.NewDigitalPin(triggerNum)
	trigger.SetDirection(embd.Out)

	// input := make(chan interface{})
	// err := echo.Watch(embd.EdgeRising, func(echo embd.DigitalPin) {
	// 	input <- true
	// })
	// if err != nil {
	// 	panic(err)
	// }

	trigger.Write(embd.High)
	time.Sleep(50 * time.Microsecond)
	trigger.Write(embd.Low)

	// startTime := time.Now() // Record time when ECHO goes high
	// <-input
	// echo.StopWatching()
	// duration := time.Since(startTime)

	duration, err := echo.TimePulse(embd.High)
	if err != nil {
		return 0, err
	}

	distance := float64(duration.Nanoseconds()) / 10000000 * 170

	echo.Close()
	trigger.Close()
	return distance, nil
}
