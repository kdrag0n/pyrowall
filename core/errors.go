package core

import (
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Check is a convenience method for panicking on errors.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func createSentryEvent(err interface{}) (event *sentry.Event) {
	var stacktrace *sentry.Stacktrace
	switch err := err.(type) {
	case error:
		stacktrace = sentry.ExtractStacktrace(err)
	}

	if stacktrace == nil {
		stacktrace = sentry.NewStacktrace()
		stacktrace.Frames = stacktrace.Frames[:len(stacktrace.Frames)-2]
	}

	event = sentry.NewEvent()
	event.Level = sentry.LevelError
	event.Exception = []sentry.Exception{{
		Value:      fmt.Sprintf("%v", err),
		Type:       reflect.TypeOf(err).String(),
		Stacktrace: stacktrace,
	}}

	return
}

// Recover will recover from a panic, log the error and stack, and report it to Sentry.
func Recover(ctx string) {
	if err := recover(); err != nil {
		var logEvent *zerolog.Event
		switch err := err.(type) {
		case error:
			logEvent = log.Error().Err(err)
		case string:
			logEvent = log.Error().Str(zerolog.ErrorFieldName, err)
		default:
			logEvent = log.Error().Fields(map[string]interface{}{zerolog.ErrorFieldName: err})
		}

		// Log error and stack
		logEvent.Stack().
			Str("type", reflect.TypeOf(err).String()).
			Str("context", ctx).
			Msg("Unhandled error")
		debug.PrintStack()

		// Report to Sentry
		hub := sentry.CurrentHub()
		client, scope := hub.Client(), hub.Scope()
		if client != nil && scope != nil {
			sentryEvent := createSentryEvent(err)
			sentryEvent.Tags["context"] = ctx
			client.CaptureEvent(sentryEvent, &sentry.EventHint{RecoveredException: err}, scope)
		}
	}
}
