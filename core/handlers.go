package core

import (
	"strings"
	"unicode"

	"github.com/kdrag0n/pyrowall/commands"
	"github.com/kdrag0n/pyrowall/util"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/rs/zerolog/log"
)

/*
 * Group 0 - command messages
 */

func (b *Bot) captionCmdHandler(eb ext.Bot, u *gotgbot.Update) error {
	u.Message.Text = u.Message.Caption
	return b.textCmdHandler(eb, u)
}

func (b *Bot) textCmdHandler(_ ext.Bot, u *gotgbot.Update) (ret error) {
	// Make sure we continue propagating events after running this handler
	ret = gotgbot.ContinueGroups{}

	// Only handle commands
	if u.Message.Text[0] != '/' {
		return
	}

	// Extract only the command invocation segment (first part) from the message text
	// This avoids using strings.Fields() to avoid excess allocations for long, whitespace-heavy messages since
	// we might not even end up invoking a command
	var cmdSeg string
	cmdEndIdx := strings.IndexFunc(u.Message.Text, unicode.IsSpace)
	if cmdEndIdx == -1 {
		cmdSeg = u.Message.Text[:util.Min(len(u.Message.Text), b.maxCmdSegLen)]
	} else {
		cmdSeg = u.Message.Text[cmdEndIdx:]
	}
	log.Debug().Str("segment", cmdSeg).Msg("Parsed command segment")

	// Handle commands directed towards specific bots with an @username suffix
	unameIdx := strings.IndexByte(cmdSeg, '@')
	if unameIdx != -1 {
		// Extract target username
		uname := cmdSeg[unameIdx:]

		// Ignore command if username doesn't match
		if uname != b.User.Username {
			return
		}
	}

	// Get and invoke command if valid
	cmdName := cmdSeg[1:]
	if cmd, ok := b.Commands[cmdName]; ok {
		// Construct context
		ctx := commands.Context{
			Update:     u,
			CmdSegment: cmdSeg,
		}

		// Call command function
		cmd.Func(ctx)
	}

	return
}

func (b *Bot) registerHandlers() {
	// Command message handlers
	b.updater.Dispatcher.AddHandlerToGroup(handlers.NewMessage(Filters.Caption, b.captionCmdHandler), 0)
	b.updater.Dispatcher.AddHandlerToGroup(handlers.NewMessage(Filters.Text, b.textCmdHandler), 0)
}
