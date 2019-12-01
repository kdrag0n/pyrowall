package core

import (
	"fmt"

	"github.com/kdrag0n/pyrowall/commands"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/rs/zerolog/log"
)

const (
	// Maximum command length accepted by Telegram.
	MaxCommandLength = 64
)

// Bot represents the state of a bot instance.
type Bot struct {
	updater *gotgbot.Updater
	Client  *ext.Bot
	User    *ext.User

	Config   *Config
	Modules  map[string]Module
	Commands map[string]commands.Command

	maxCmdSegLen int
}

// NewBot returns a new Bot instance with the provided Config.
func NewBot(config *Config) *Bot {
	return &Bot{
		updater: nil,
		Client:  nil,

		Config:   config,
		Modules:  make(map[string]Module),
		Commands: make(map[string]commands.Command),
	}
}

func (b *Bot) startWebhooks() error {
	webhook := gotgbot.Webhook{
		Serve:          b.Config.Telegram.WebhookAddress,
		ServePort:      b.Config.Telegram.WebhookPort,
		ServePath:      b.Config.Telegram.WebhookPath,
		URL:            b.Config.Telegram.WebhookURL,
		MaxConnections: b.Config.Telegram.WebhookConnectionLimit,
	}

	b.updater.StartWebhook(webhook)
	ok, err := b.updater.SetWebhook(b.updater.Bot.Token, webhook)
	if err != nil {
		return fmt.Errorf("start webhooks: %w", err)
	}

	if !ok {
		return fmt.Errorf("set webhook: %w", err)
	}

	return nil
}

func (b *Bot) startPolling() error {
	err := b.updater.StartPolling()
	if err != nil {
		return fmt.Errorf("start polling: %w", err)
	}

	return nil
}

func (b *Bot) startUpdater() error {
	if b.Config.Telegram.UseWebhooks {
		return b.startWebhooks()
	} else {
		return b.startPolling()
	}
}

// Start connects to Telegram and starts the bot.
func (b *Bot) Start() error {
	// Load modules
	err := b.LoadModules()
	if err != nil {
		return err
	}

	// Create updater
	b.updater, err = gotgbot.NewUpdater(b.Config.Telegram.Token)
	if err != nil {
		return fmt.Errorf("create updater: %w", err)
	}

	// Set client and get self user
	b.Client = b.updater.Dispatcher.Bot
	b.User, err = b.Client.GetMe()
	if err != nil {
		return fmt.Errorf("fetch user info: %w", err)
	}
	log.Info().Interface("user", b.User).Msg("Fetched self user data")

	// Calculate maximum command segment length (optimization for large messages)
	// Prefix length (/) + max command name length + mention prefix length (@) + username length
	b.maxCmdSegLen = 1 + MaxCommandLength + 1 + len(b.User.Username)
	log.Debug().Int("len", b.maxCmdSegLen).Msg("Calculated max command segment length")

	// Register handlers
	b.registerHandlers()

	// Start updater
	return b.startUpdater()
}
