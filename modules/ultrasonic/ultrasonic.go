package ultrasonic

import (
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
	data      = make(map[string]conf)
	status    chan int
	sensor    string
	lastVal   = make(map[string]float64)
	messures  = make(map[int]struct {
		distance float64
		run      bool
		value    int
	})
)

func Register() {
	m := modules.Module{name, ultrasonicModule{}, true}
	m.Register()
}

func (m ultrasonicModule) Config() interface{} {
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
	for sensor, value := range data {
		echos[sensor], _ = embd.NewDigitalPin(value.Echo)
		echos[sensor].SetDirection(embd.In)
		triggers[sensor], _ = embd.NewDigitalPin(value.Trigger)
		triggers[sensor].SetDirection(embd.Out)

		status = make(chan int)
		err := echos[sensor].Watch(embd.EdgeBoth, func(echo embd.DigitalPin) {
			read, _ := echo.Read()
			status <- read
		})
		if err != nil {
			panic(err)
		}
		e := event.NewEvent(name, sensor, "web")
		e.AddData("start", "1")
		e.SendEventToWeb()
	}
	go measure(status)

	for running {
		for sensor = range data {
			triggers[sensor].Write(embd.High)
			time.Sleep(1 * time.Millisecond)
			triggers[sensor].Write(embd.Low)

			time.Sleep(70 * time.Millisecond)
			status <- timeout
			time.Sleep(5 * time.Millisecond)
		}
	}

	for key := range data {
		echos[key].StopWatching()
		echos[key].Close()
		triggers[key].Close()
		e := event.NewEvent(name, key, "web")
		e.AddData("value", "stop")
		e.SendEventToWeb()
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
		if lastVal[sensor] > 0 {
			distance = (lastVal[sensor]*2 + distance) / 3
		}
		lastVal[sensor] = distance

		value := strconv.FormatFloat(distance, 'f', 2, 64)
		e := event.NewEvent(name, sensor, "distance")
		e.AddData("value", value)
		e.AddData("posDegree", strconv.Itoa(data[sensor].Position.Degree))
		e.AddData("posDistance", strconv.Itoa(data[sensor].Position.Distance))
		go e.SendEventToAll()
	}
}
