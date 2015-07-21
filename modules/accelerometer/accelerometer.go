package accelerometer

import (
	"gotank/components/mpu6050"
	"gotank/event"
	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
	"gotank/modules"
	"strconv"
	"time"
)

const name = "accelerometer"

type accelerometerModule struct{}
type conf struct {
	Degree, Min, Rotation int
}

var (
	running       = false
	data          map[string]conf
	accelerometer map[string]*mpu6050.MPU6050Driver
)

func Register() {
	data = make(map[string]conf)
	m := modules.Module{name, accelerometerModule{}, true}
	m.Register()
}

func (m accelerometerModule) Config() interface{} {
	return &data
}

func (m accelerometerModule) Start() {
	running = true
	go read()
}

func (m accelerometerModule) Stop() {
	running = false
}

func (m accelerometerModule) GetEvent(e event.Event) {
}

func (m accelerometerModule) Active() []string {
	var active []string
	if data == nil {
		return active
	}

	for key := range data {
		active = append(active, key)
	}
	return active
}

func read() {
	accelerometer = make(map[string]*mpu6050.MPU6050Driver)
	for key := range data {
		accelerometer[key] = mpu6050.NewMPU6050Driver(embd.NewI2CBus(1))
	}

	for running {
		for key := range data {
			e := event.NewEvent(name, key, "accelerometer")

			temperature := accelerometer[key].GetTemperature()
			e.AddData("temperature", strconv.FormatFloat(temperature, 'f', 2, 64))

			accX, accY, accZ := accelerometer[key].GetAccelerometer()
			e.AddData("accX", string(accX))
			e.AddData("accY", string(accY))
			e.AddData("accZ", string(accZ))

			gyroX, gyroY, gyroZ := accelerometer[key].GetGyroscope()
			e.AddData("gyroX", string(gyroX))
			e.AddData("gyroY", string(gyroY))
			e.AddData("gyroZ", string(gyroZ))

			e.SendEventToAll()
		}
		time.Sleep(50 * time.Millisecond)
	}
}
