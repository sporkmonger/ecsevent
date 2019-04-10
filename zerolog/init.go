package zerolog

import (
	"time"

	"github.com/rs/zerolog"

	"github.com/sporkmonger/ecslog"
)

func init() {
	zerologFieldInit()
}

func zerologFieldInit() {
	zerolog.LevelFieldName = ecslog.FieldLogLevel
	zerolog.MessageFieldName = ecslog.FieldMessage
	zerolog.ErrorFieldName = ecslog.FieldErrorMessage
	zerolog.TimestampFieldName = ecslog.FieldTimestamp
	zerolog.TimeFieldFormat = time.RFC3339Nano
	// Do not use caller for now.
}
