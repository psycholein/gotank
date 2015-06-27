package motorshield

import (
	"reflect"
	"strings"
)

type control struct {
	key   string
	value string
}

func (c control) handle() {
	_, l := left[c.key]
	_, r := right[c.key]
	if len(c.value) < 2 || !l || !r {
		return
	}

	methodName := strings.ToUpper(c.value[0:1]) + c.value[1:len(c.value)]
	method := reflect.ValueOf(&c).MethodByName(methodName)
	if method.IsValid() {
		method.Call([]reflect.Value{})
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
