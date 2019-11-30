package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/kdrag0n/pyrowall/core"

	"github.com/rs/zerolog/log"
)

func startPprof(config *core.Config) {
	// Construct listen address
	listenAddr := fmt.Sprintf("%s:%d", config.Pprof.ServerAddress, config.Pprof.ServerPort)
	log.Info().Str("address", listenAddr).Msg("Starting pprof server")

	// Swap ServeMux to isolate server
	pprofMux := http.DefaultServeMux
	http.DefaultServeMux = http.NewServeMux()

	// Start server
	go func() {
		err := http.ListenAndServe(listenAddr, pprofMux)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start pprof server")
		}
	}()
}
