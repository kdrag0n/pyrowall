package core

import (
	"fmt"
	"time"

	"github.com/kdrag0n/pyrowall/commands"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/go-pg/pg/v9"
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

	DB *pg.DB

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
	ext.DefaultTgBotGetter.Client.Timeout = time.Second * 3

	if b.Config.Telegram.UseWebhooks {
		return b.startWebhooks()
	} else {
		return b.startPolling()
	}
}

func (b *Bot) fillUserInfo() (err error) {
	// Get user data
	b.User, err = b.Client.GetMe()
	if err != nil {
		return fmt.Errorf("fetch user info: %w", err)
	}
	log.Info().Interface("user", b.User).Msg("Fetched self user data")

	// Calculate maximum command segment length (optimization for large messages)
	// Prefix length (/) + max command name length + mention prefix length (@) + username length
	b.maxCmdSegLen = 1 + MaxCommandLength + 1 + len(b.User.Username)
	log.Debug().Int("len", b.maxCmdSegLen).Msg("Calculated max command segment length")

	return
}

// Start initiates the core network connections and starts the bot.
func (b *Bot) Start() (err error) {
	// Connect to database
	log.Info().Msg("Connecting to database...")
	err = b.connectToDB()
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	// Create updater
	log.Info().Msg("Connecting to Telegram...")
	b.updater, err = gotgbot.NewUpdater(b.Config.Telegram.Token)
	if err != nil {
		return fmt.Errorf("create updater: %w", err)
	}

	// Set client and fill self user info
	b.Client = b.updater.Dispatcher.Bot
	log.Info().Msg("Fetching user info...")
	err = b.fillUserInfo()
	if err != nil {
		return
	}

	// Register handlers
	b.registerHandlers()

	// Load modules
	err = b.LoadModules()
	if err != nil {
		return
	}

	// Start updater
	return b.startUpdater()
}
