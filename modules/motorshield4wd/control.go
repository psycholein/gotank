package motorshield4wd

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
	_, l1 := left1[c.key]
	_, r1 := right1[c.key]
	_, l2 := left2[c.key]
	_, r2 := right2[c.key]
	value := c.data["value"]
	if len(value) < 2 || !l1 || !r1 || !l2 || !r2 {
		return
	}

	methodName := strings.ToUpper(value[0:1]) + value[1:len(value)]
	method := reflect.ValueOf(&c).MethodByName(methodName)
	if method.IsValid() {
		method.Call([]reflect.Value{})
		e := event.NewEvent(name, c.key, "control_feeback")
		e.AddData("value", value)
		go e.SendEventToAll()
	}
}

func (c control) Forward() {
	left1[c.key].Forward()
	left2[c.key].Forward()
	right1[c.key].Forward()
	right2[c.key].Forward()
}

func (c control) Backward() {
	left1[c.key].Backward()
	left2[c.key].Backward()
	right1[c.key].Backward()
	right2[c.key].Backward()
}

func (c control) Turnleft() {
	left1[c.key].Backward()
	left2[c.key].Backward()
	right1[c.key].Forward()
	right2[c.key].Forward()
}

func (c control) Turnright() {
	left1[c.key].Forward()
	left2[c.key].Forward()
	right1[c.key].Backward()
	right2[c.key].Backward()
}

func (c control) Stop() {
	left1[c.key].Stop()
	left2[c.key].Stop()
	right1[c.key].Stop()
	right2[c.key].Stop()
}

func (c control) Speed(value int) {
	left1[c.key].Speed(value)
	left2[c.key].Speed(value)
	right1[c.key].Speed(value)
	right2[c.key].Speed(value)
}
