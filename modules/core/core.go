package core

import (
	"github.com/kdrag0n/pyrowall/commands"
	"github.com/kdrag0n/pyrowall/core"
)

// Module contains the state for an instance of this module.
type Module struct {
	Bot *core.Bot
}

// Info returns basic information about this module.
func (m *Module) Info() core.ModuleInfo {
	return core.ModuleInfo{
		Name: "Core",
	}
}

// Commands returns a list of commands provided by this module.
func (m *Module) Commands() []commands.Command {
	return []commands.Command{
		{
			Name:        "start",
			Description: "Start an interaction with me.",
			Func:        m.cmdStart,
		},
	}
}

// NewModule returns a new instance of this module.
func NewModule(bot *core.Bot) (core.Module, error) {
	return &Module{
		Bot: bot,
	}, nil
}

func init() {
	core.RegisterModule("Core", NewModule)
}

/*
 * Commands
 */

func (m *Module) cmdStart(c commands.Context) {
	_, err := c.Message.ReplyText("Hello!")
	core.Check(err)
}
