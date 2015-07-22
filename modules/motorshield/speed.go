package motorshield

import (
	"gotank/event"
	"strconv"
)

func handleSpeed(key string, data map[string]string) {
	_, l := left[key]
	_, r := right[key]
	value, _ := strconv.Atoi(data["value"])
	if value <= 0 || !l || !r {
		return
	}
	if value > 1024 {
		value = 1024
	}

	setSpeed(key, value)

	e := event.NewEvent(name, key, "speed")
	e.AddData("value", data["value"])
	go e.SendEventToAll()
}

func setSpeed(key string, value int) {
	left[key].Speed(value)
	right[key].Speed(value)
}
