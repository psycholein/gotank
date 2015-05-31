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
type position struct {
	Distance, Degree int
}
type conf struct {
	Trigger, Echo int
	Position      position
}

var (
	running   = false
	startTime = time.Now()
	data      map[string]conf
)

func Register() {
	data = make(map[string]conf)
	m := modules.Module{name, ultrasonicModule{}}
	m.Register(data)
}

func (m ultrasonicModule) Start() {
	fmt.Println(data)
	running = true
	go distance()
}

func (m ultrasonicModule) Stop() {
	running = false
}

func distance() {
	echo, _ := embd.NewDigitalPin(17)
	echo.SetDirection(embd.In)
	trigger, _ := embd.NewDigitalPin(4)
	trigger.SetDirection(embd.Out)

	input := make(chan int)
	err := echo.Watch(embd.EdgeBoth, func(echo embd.DigitalPin) {
		read, _ := echo.Read()
		input <- read
	})
	if err != nil {
		panic(err)
	}
	go measure(input)

	for running {
		trigger.Write(embd.High)
		time.Sleep(50 * time.Microsecond)
		trigger.Write(embd.Low)

		time.Sleep(50 * time.Millisecond)
	}

	echo.StopWatching()
	echo.Close()
	trigger.Close()
	close(input)
}

func measure(c chan int) {
	for val := range c {
		if val == 1 {
			startTime = time.Now()
			continue
		}
		duration := time.Since(startTime)
		distance := float64(duration.Nanoseconds()) / 10000000 * 171.5
		fmt.Println(distance)
	}
}
