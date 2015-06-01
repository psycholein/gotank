package ultrasonic

import (
	"fmt"
	"gotank/modules"
	"time"

	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
)

const name = "ultrasonic"
const (
	measureEnd   = 0
	measureStart = 1
	timeout      = 2
)

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
	status    chan int
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

	status = make(chan int, 1)
	err := echo.Watch(embd.EdgeBoth, func(echo embd.DigitalPin) {
		read, _ := echo.Read()
		status <- read
	})
	if err != nil {
		panic(err)
	}
	go measure(status)

	for running {
		trigger.Write(embd.High)
		time.Sleep(50 * time.Microsecond)
		trigger.Write(embd.Low)

		time.Sleep(50 * time.Millisecond)
		status <- timeout
		time.Sleep(1 * time.Millisecond)
	}

	echo.StopWatching()
	echo.Close()
	trigger.Close()
	close(status)
}

func measure(status chan int) {
	run := false
	for val := range status {
		if val == timeout {
			run = false
			continue
		}
		if val == measureStart {
			run = true
			startTime = time.Now()
			continue
		}
		if !run {
			fmt.Println("Distance timeout")
			continue
		}
		duration := time.Since(startTime)
		distance := float64(duration.Nanoseconds()) / 10000000 * 171.5
		fmt.Println(distance)
	}
}
