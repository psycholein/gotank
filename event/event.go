package event

type Event struct {
	Module string
	Name   string
	Task   string
	Data   map[string]string
}

type eventType struct {
	event Event
	web   bool
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

func InitEvents() {
	register = make(map[string][]EventFunc)
	eventChan = make(chan eventType)
	go handleEvents()
}

func SendEvent(e Event) {
	eventChan <- eventType{e, false}
}

func SendWebEvent(e Event) {
	eventChan <- eventType{e, true}
}

func RegisterEvent(srcModule string, destModule EventFunc) {
	register[srcModule] = append(register[srcModule], destModule)
}

func handleEvents() {
	for ec := range eventChan {
		if items, ok := register[ec.event.Module]; ok {
			for _, item := range items {
				item(ec.event)
			}
		}
		if !ec.web {
			if items, ok := register["_all"]; ok {
				for _, item := range items {
					item(ec.event)
				}
			}
		}
	}
}

func Stop() {
	close(eventChan)
}
