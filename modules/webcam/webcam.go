package webcam

import (
	"fmt"
	"gotank/event"
	"gotank/modules"
	"os/exec"
)

const name = "webcam"

type webcamModule struct{}
type conf struct {
	Device string
}

var (
	running = false
	data    = make(map[string]conf)
)

func Register() {
	m := modules.Module{name, webcamModule{}, true}
	m.Register()
}

func (m webcamModule) Config() interface{} {
	return &data
}

func (m webcamModule) Start() {
	running = true
	for key := range data {
		e := event.NewEvent(name, key, "web")
		e.AddData("start", "1")
		e.SendEventToWeb()
	}
	go mjpgStreamer()
}

func (m webcamModule) Stop() {
	running = false
	for key := range data {
		e := event.NewEvent(name, key, "web")
		e.AddData("stop", "1")
		e.SendEventToWeb()
	}
}

func (m webcamModule) GetEvent(e event.Event) {
}

func (m webcamModule) Active() []string {
	var active []string
	if data == nil {
		return active
	}

	for key := range data {
		active = append(active, key)
	}
	return active
}

func mjpgStreamer() {
	cmd := "mjpg_streamer"
	input := "input_uvc.so -d /dev/video0 -l off -r 320x240 -f 25"
	output := "output_http.so -p 8080"
	exec.Command(cmd, "-i", input, "-o", output).Run()
	fmt.Println("webcam call")
}
