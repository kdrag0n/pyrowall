package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
)

// LogrusInterposer translates Logrus calls to zerolog.
type LogrusInterposer struct{}

func (f *LogrusInterposer) levelToEvent(level logrus.Level) *zerolog.Event {
	switch level {
	case logrus.TraceLevel:
		return log.Trace()
	case logrus.DebugLevel:
		return log.Debug()
	case logrus.InfoLevel:
		return log.Info()
	case logrus.WarnLevel:
		return log.Warn()
	case logrus.ErrorLevel:
		return log.Error()
	case logrus.FatalLevel:
		return log.Fatal()
	case logrus.PanicLevel:
		return log.Panic()
	default:
		return log.Info()
	}
}

// Format passes a log entry to zerolog.
func (f *LogrusInterposer) Format(entry *logrus.Entry) (_ []byte, _ error) {
	// Create zerolog event with the correct level
	zlEvent := f.levelToEvent(entry.Level)

	// Add fields
	zlEvent.Fields(entry.Data)

	// Add caller info if present
	if entry.HasCaller() {
		zlEvent.Str(zerolog.CallerFieldName, zerolog.CallerMarshalFunc(entry.Caller.File, entry.Caller.Line))

		if entry.Caller.Function != "" {
			zlEvent.Str(logrus.FieldKeyFunc, entry.Caller.Function)
		}
	}

	// Log with message and return
	zlEvent.Msg(entry.Message)
	return
}
