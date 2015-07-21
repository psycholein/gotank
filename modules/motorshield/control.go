package motorshield

import (
	"gotank/event"
	"reflect"
	"strings"
)

type control struct {
	key  string
	data map[string]string
}

func (c control) handleControl() {
	_, l := left[c.key]
	_, r := right[c.key]
	value := c.data["value"]
	if len(value) < 2 || !l || !r {
		return
	}

	methodName := strings.ToUpper(value[0:1]) + value[1:len(value)]
	method := reflect.ValueOf(&c).MethodByName(methodName)
	if method.IsValid() {
		method.Call([]reflect.Value{})
		e := event.NewEvent(name, c.key, "control")
		e.AddData("value", value)
		e.SendEventToAll()
	}
}

func (c control) Forward() {
	left[c.key].Forward()
	right[c.key].Forward()
}

func (c control) Backward() {
	left[c.key].Backward()
	right[c.key].Backward()
}

func (c control) Turnleft() {
	left[c.key].Backward()
	right[c.key].Forward()
}

func (c control) Turnright() {
	left[c.key].Forward()
	right[c.key].Backward()
}

func (c control) Stop() {
	left[c.key].Stop()
	right[c.key].Stop()
}

func (c control) Speed(value int) {
	left[c.key].Speed(value)
	right[c.key].Speed(value)
}
