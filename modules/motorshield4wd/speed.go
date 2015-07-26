package motorshield4wd

import (
	"gotank/event"
	"strconv"
)

func handleSpeed(key string, data map[string]string) {
	_, l1 := left1[key]
	_, r1 := right1[key]
	_, l2 := left2[key]
	_, r2 := right2[key]
	value, _ := strconv.Atoi(data["value"])
	if value <= 0 || !l1 || !r1 || !l2 || !r2 {
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
	left1[key].Speed(value)
	left2[key].Speed(value)
	right1[key].Speed(value)
	right2[key].Speed(value)
}
