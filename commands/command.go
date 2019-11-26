package commands

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
