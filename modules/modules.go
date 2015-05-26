package modules

//_ "gotank/libs/embd"
//_ "gotank/libs/embd/host/all" // select the right board

type event struct {
	Name string
	Data string
}

type moduleInterface interface {
	Start(config interface{})
	Stop()
}

type module struct {
	Name string
	Attr moduleInterface
}

var availableModules map[string]module

func (m module) register() {
	if availableModules == nil {
		availableModules = make(map[string]module)
	}
	availableModules[m.Name] = m
}

type mainConfig struct {
	Modules []string
}

func InitModule(name string, config interface{}) {
	if val, ok := availableModules[name]; ok {
		val.Attr.Start(config)
		return
	}
}
