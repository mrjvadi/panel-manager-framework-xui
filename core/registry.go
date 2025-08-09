package core

var registry = map[string]DriverFactory{}

func Register(name string, f DriverFactory) { registry[name] = f }
func Factory(name string) (DriverFactory, bool) { f, ok := registry[name]; return f, ok }
