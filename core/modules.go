package core

import (
	"fmt"

	"github.com/kdrag0n/pyrowall/commands"
)

// ModuleInfo contains the information of a Module.
type ModuleInfo struct {
	Name string
}

// Module defines a bot module.
type Module interface {
	Info() ModuleInfo
	Commands() []commands.Command
}

// A ModuleConstructor returns an instance of its module tied to the given Bot.
type ModuleConstructor func(*Bot) (Module, error)

// Modules contains the name and constructor of each registered Module.
var Modules = make(map[string]ModuleConstructor)

// RegisterModule registers a module with the given name and constructor.
func RegisterModule(name string, constructor ModuleConstructor) {
	if _, ok := Modules[name]; ok {
		panic(fmt.Errorf("attempted to register module under occupied name '%s'", name))
	}

	Modules[name] = constructor
}
