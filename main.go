package main

import (
	"os"
	"runtime"

	"github.com/kdrag0n/pyrowall/core"
	_ "github.com/kdrag0n/pyrowall/modules"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogging() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "Jan 02 15:04:05"})
}

func readConfig() (config *core.Config) {
	log.Info().Msg("Reading config...")
	cfgName := "config.toml"
	if len(os.Args) > 1 {
		cfgName = os.Args[1]
	}

	config, err := core.ReadConfigFile(cfgName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}

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
	setupLogging()

	config := readConfig()
	startBot(config)

	log.Info().Msg("Bot started")
	runtime.Goexit()
}
