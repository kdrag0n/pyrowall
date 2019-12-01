package commands

import "github.com/PaulSonOfLars/gotgbot"

// CommandFunc represents a command function that takes no message arguments.
type CommandFunc = func(Context)

// Command describes a bot command.
type Command struct {
	Name        string
	Description string
	Usage       string
	Aliases     []string
	Func        CommandFunc
}

// Invoke invokes a Command with the given arguments.
func (cmd *Command) Invoke(update *gotgbot.Update, cmdSeg string) {
	// Construct context
	ctx := newContext(update, cmdSeg)
	cmd.Func(ctx)
}
