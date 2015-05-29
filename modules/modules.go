package modules

//_ "gotank/libs/embd"
//_ "gotank/libs/embd/host/all" // select the right board

type moduleInterface interface {
	Start(config interface{})
	Stop()
}

type Module struct {
	Name string
	Attr moduleInterface
}

var availableModules map[string]Module

func (m Module) Register() {
	if availableModules == nil {
		availableModules = make(map[string]Module)
	}
	availableModules[m.Name] = m
}

func InitModule(name string, config interface{}) {
	if val, ok := availableModules[name]; ok {
		val.Attr.Start(config)
		return
	}
}
