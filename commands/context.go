package commands

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

// Context contains the context information used to invoke a command.
type Context struct {
	Update *gotgbot.Update

	// Convenience fields
	Message *ext.Message
	Chat    *ext.Chat
	User    *ext.User

	// Miscellaneous
	CmdSegment string

	// Lazy values
	args        []string
	haveRawArgs bool
	rawArgs     string
}

func newContext(update *gotgbot.Update, cmdSeg string) Context {
	return Context{
		Update: update,

		Message: update.EffectiveMessage,
		Chat:    update.EffectiveChat,
		User:    update.EffectiveUser,

		CmdSegment: cmdSeg,
	}
}

// Args returns a slice of whitespace-separated arguments from the command message.
func (c *Context) Args() []string {
	if c.args == nil {
		c.args = strings.Fields(c.Message.Text)[1:]
	}

	return c.args
}

// RawArgs returns a string with everything in the command message except the command invocation segment.
func (c *Context) RawArgs() string {
	if !c.haveRawArgs {
		c.rawArgs = c.Message.Text[len(c.CmdSegment):]
		c.haveRawArgs = true
	}

	return c.rawArgs
}
