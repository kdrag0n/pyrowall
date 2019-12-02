package main

import (
	"os"
	"runtime"

	"github.com/kdrag0n/pyrowall/util"

	"github.com/getsentry/sentry-go"

	"github.com/sirupsen/logrus"

	"github.com/kdrag0n/pyrowall/core"
	_ "github.com/kdrag0n/pyrowall/modules"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogging(config *core.Config) {
	// Configure zerolog
	if config.Logging.Enable {
		level, err := zerolog.ParseLevel(config.Logging.Level)
		util.PanicIf(err)
		zerolog.SetGlobalLevel(level)

		if config.Logging.Format == "console" {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "Jan 02 15:04:05"})
		}
	} else {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	// Register logrus->zerolog interposer and let zerolog handle levels
	logrus.SetFormatter(&LogrusInterposer{})
	logrus.SetLevel(logrus.TraceLevel)
}

func setupSentry(config *core.Config) {
	log.Info().Msg("Initializing Sentry error reporting...")

	// Default to Git commit hash for release
	if config.Sentry.Release == "" {
		config.Sentry.Release = GitCommit
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         config.Sentry.DSN,
		ServerName:  config.Sentry.ServerName,
		Release:     config.Sentry.Release,
		Environment: config.Sentry.Environment,

		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Sentry client")
	}
}

func readConfig() (config *core.Config) {
	cfgName := "config.toml"
	if len(os.Args) > 1 {
		cfgName = os.Args[1]
	}

	config, err := core.ReadConfigFile(cfgName)
	util.PanicIf(err)

	return
}

func startBot(config *core.Config) *core.Bot {
	log.Info().Str("commit", GitCommit).Msg("Starting bot...")
	bot := core.NewBot(config)

	err := bot.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start bot")
	}

	return bot
}

func main() {
	config := readConfig()
	setupLogging(config)
	setupSentry(config)

	if config.Pprof.EnableServer {
		startPprof(config)
	}

	startBot(config)

	log.Info().Msg("Bot started")
	runtime.Goexit()
}
