package core

import (
	"strings"
	"unicode"

	"github.com/kdrag0n/pyrowall/util"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/rs/zerolog/log"
)

/*
 * Group -20: leave channels
 */

func channelPredicate(m *ext.Message) bool {
	return m.Chat.Type == "channel"
}

func (b *Bot) channelLeaveHandler(_ ext.Bot, u *gotgbot.Update) (ret error) {
	log.Debug().Str("name", u.EffectiveChat.Title).Int("id", u.EffectiveChat.Id).Msg("Leaving channel")

	// Leave immediately since we already checked the chat type in the filter
	_, err := b.Client.LeaveChat(u.EffectiveChat.Id)
	if err != nil {
		log.Err(err).Msg("Failed to leave channel")
	}

	return
}

/*
 * Group 0: command messages
 */

// Run first & propagate to allow commands to be used in photo captions
func (b *Bot) captionCmdHandler(eb ext.Bot, u *gotgbot.Update) error {
	u.Message.Text = u.Message.Caption
	return gotgbot.ContinueGroups{}
}

func textCmdPredicate(m *ext.Message) bool {
	return m.Text != "" && m.Text[0] == '/'
}

func (b *Bot) textCmdHandler(_ ext.Bot, u *gotgbot.Update) (ret error) {
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
		// Invoke command
		cmd.Invoke(u, cmdSeg)
	}

	return
}

func (b *Bot) registerHandlers() {
	dsp := b.updater.Dispatcher

	// Channel leave handler
	channelHandler := handlers.NewMessage(channelPredicate, b.channelLeaveHandler)
	channelHandler.AllowChannel = true
	dsp.AddHandlerToGroup(channelHandler, -20)

	// Command message handlers
	dsp.AddHandlerToGroup(handlers.NewMessage(Filters.Caption, b.captionCmdHandler), 0)
	dsp.AddHandlerToGroup(handlers.NewMessage(textCmdPredicate, b.textCmdHandler), 0)
}
