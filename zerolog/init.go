package zerolog

import (
	"time"

	"github.com/rs/zerolog"

	"github.com/sporkmonger/slog"
)

func init() {
	zerologFieldInit()
}

func zerologFieldInit() {
	zerolog.LevelFieldName = slog.FieldLogLevel
	zerolog.MessageFieldName = slog.FieldMessage
	zerolog.ErrorFieldName = slog.FieldErrorMessage
	zerolog.TimestampFieldName = slog.FieldTimestamp
	zerolog.TimeFieldFormat = time.RFC3339Nano
	// Do not use caller for now.
}
