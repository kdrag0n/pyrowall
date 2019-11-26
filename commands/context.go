package commands

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot"
)

// Context contains the context information used to invoke a command.
type Context struct {
	Update *gotgbot.Update

	// Miscellaneous
	CmdSegment string

	// Lazy values
	args        []string
	haveRawArgs bool
	rawArgs     string
}

// Args returns a slice of whitespace-separated arguments from the command message.
func (c *Context) Args() []string {
	if c.args == nil {
		c.args = strings.Fields(c.Update.Message.Text)[1:]
	}

	return c.args
}

// RawArgs returns a string with everything in the command message except the command invocation segment.
func (c *Context) RawArgs() string {
	if !c.haveRawArgs {
		c.rawArgs = c.Update.Message.Text[len(c.CmdSegment):]
		c.haveRawArgs = true
	}

	return c.rawArgs
}
