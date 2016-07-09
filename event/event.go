package event

const (
	all   = 0
	local = 1
	web   = 2
)

type Event struct {
	Module string
	Name   string
	Task   string
	Data   map[string]string
}

type eventType struct {
	event    Event
	location int
}

type EventFunc func(Event)

var eventChan chan eventType
var register map[string][]EventFunc

func NewEvent(module, name, task string) (e Event) {
	e = Event{module, name, task, make(map[string]string)}
	return
}

func (e Event) AddData(name string, value string) {
	e.Data[name] = value
}

func (e Event) SendEventToAll() {
	eventChan <- eventType{e, all}
}

func (e Event) SendEventToLocal() {
	eventChan <- eventType{e, local}
}

func (e Event) SendEventToWeb() {
	eventChan <- eventType{e, web}
}

func InitEvents() {
	register = make(map[string][]EventFunc)
	eventChan = make(chan eventType)
	go handleEvents()
}

func RegisterEvent(srcModule string, destModule EventFunc) {
	register[srcModule] = append(register[srcModule], destModule)
}

func handleEvents() {
	for ec := range eventChan {
		if ec.location == all || ec.location == local {
			if items, ok := register[ec.event.Module]; ok {
				for _, item := range items {
					go item(ec.event)
				}
			}
		}
		if ec.location == all || ec.location == web {
			if items, ok := register["_web"]; ok {
				for _, item := range items {
					go item(ec.event)
				}
			}
		}
		if items, ok := register["_all"]; ok {
			for _, item := range items {
				go item(ec.event)
			}
		}
	}
}

func Stop() {
	close(eventChan)
}
