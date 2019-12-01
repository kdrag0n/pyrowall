package main

import (
	"os"
	"runtime"

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
		core.Check(err)
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

func readConfig() (config *core.Config) {
	cfgName := "config.toml"
	if len(os.Args) > 1 {
		cfgName = os.Args[1]
	}

	config, err := core.ReadConfigFile(cfgName)
	core.Check(err)

	return
}

func startBot(config *core.Config) *core.Bot {
	log.Info().Msg("Starting bot...")
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

	if config.Pprof.EnableServer {
		startPprof(config)
	}

	startBot(config)

	log.Info().Msg("Bot started")
	runtime.Goexit()
}
