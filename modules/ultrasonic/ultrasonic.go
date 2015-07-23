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
type messureType struct {
	distance  float64
	run       bool
	name      string
	startTime time.Time
}
type statusType struct {
	pin   int
	value int
}

var (
	running  = false
	data     = make(map[string]conf)
	status   = make(chan statusType)
	messures = make(map[int]messureType)
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
		messures[value.Echo] = messureType{name: sensor}

		err := echos[sensor].Watch(embd.EdgeBoth, func(echo embd.DigitalPin) {
			read, _ := echo.Read()
			status <- statusType{echo.N(), read}
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
		for sensor, val := range data {
			triggers[sensor].Write(embd.High)
			time.Sleep(1 * time.Millisecond)
			triggers[sensor].Write(embd.Low)

			time.Sleep(70 * time.Millisecond)
			status <- statusType{val.Echo, timeout}
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

func measure(status chan statusType) {
	for val := range status {
		sensor := messures[val.pin]
		if val.value == timeout {
			sensor.run = false
			messures[val.pin] = sensor
			continue
		}
		if val.value == measureStart {
			sensor.run = true
			sensor.startTime = time.Now()
			messures[val.pin] = sensor
			continue
		}
		if !sensor.run {
			continue
		}

		duration := time.Since(sensor.startTime)
		distance := float64(duration.Nanoseconds()) / 10000000 * 171.5
		if sensor.distance > 0 {
			distance = (sensor.distance*2 + distance) / 3
		}
		sensor.distance = distance
		messures[val.pin] = sensor

		value := strconv.FormatFloat(distance, 'f', 2, 64)
		e := event.NewEvent(name, sensor.name, "distance")
		e.AddData("value", value)
		e.AddData("posDegree", strconv.Itoa(data[sensor.name].Position.Degree))
		e.AddData("posDistance", strconv.Itoa(data[sensor.name].Position.Distance))
		go e.SendEventToAll()
	}
}
