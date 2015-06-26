package event

type Event struct {
	Module string
	Name   string
	Task   string
	Value  string
}

type eventType struct {
	event Event
	web   bool
}

type EventFunc func(Event)

var event chan eventType
var register map[string][]EventFunc

func InitEvents() {
	register = make(map[string][]EventFunc)
	event = make(chan eventType)
	go handleEvents()
}

func SendEvent(e Event) {
	event <- eventType{e, false}
}

func SendWebEvent(e Event) {
	event <- eventType{e, true}
}

func RegisterEvent(srcModule string, destModule EventFunc) {
	register[srcModule] = append(register[srcModule], destModule)
}

func handleEvents() {
	for e := range event {
		if items, ok := register[e.event.Module]; ok {
			for _, item := range items {
				item(e.event)
			}
		}
		if !e.web {
			if items, ok := register["_all"]; ok {
				for _, item := range items {
					item(e.event)
				}
			}
		}
	}
}

func Stop() {
	close(event)
}
