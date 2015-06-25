package ultrasonic

import (
	"fmt"
	"gotank/event"
	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
	"gotank/modules"
	"strconv"
	"time"
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
	sensor    string
)

func Register() {
	m := modules.Module{name, ultrasonicModule{}, true}
	m.Register()
}

func (m ultrasonicModule) Config() interface{} {
	data = make(map[string]conf)
	return &data
}

func (m ultrasonicModule) Start() {
	running = true
	go distance()
}

func (m ultrasonicModule) Stop() {
	running = false
}

func (m ultrasonicModule) GetEvent(e event.Event) {
	fmt.Println("ultrasonic", e)
}

func (m ultrasonicModule) Active() []string {
	var active []string
	if data == nil {
		return active
	}

	for key := range data {
		active = append(active, key)
	}
	return active
}

func distance() {
	triggers := make(map[string]embd.DigitalPin)
	echos := make(map[string]embd.DigitalPin)
	for key, value := range data {
		echos[key], _ = embd.NewDigitalPin(value.Echo)
		echos[key].SetDirection(embd.In)
		triggers[key], _ = embd.NewDigitalPin(value.Trigger)
		triggers[key].SetDirection(embd.Out)

		status = make(chan int, 1)
		err := echos[key].Watch(embd.EdgeBoth, func(echo embd.DigitalPin) {
			read, _ := echo.Read()
			status <- read
		})
		if err != nil {
			panic(err)
		}
		event.SendEvent(event.Event{name, key, "web", "start"})
	}
	go measure(status)

	for running {
		for sensor = range data {
			triggers[sensor].Write(embd.High)
			time.Sleep(50 * time.Microsecond)
			triggers[sensor].Write(embd.Low)

			time.Sleep(50 * time.Millisecond)
			status <- timeout
			time.Sleep(1 * time.Millisecond)
		}
	}

	for key := range data {
		echos[key].StopWatching()
		echos[key].Close()
		triggers[key].Close()
		event.SendEvent(event.Event{name, key, "web", "stop"})
	}
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
			continue
		}
		duration := time.Since(startTime)
		distance := float64(duration.Nanoseconds()) / 10000000 * 171.5
		event.SendEvent(event.Event{name, sensor, "distance", strconv.FormatFloat(distance, 'f', 2, 64)})
	}
}
