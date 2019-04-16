package zerolog

import (
	"time"

	"github.com/rs/zerolog"

	"github.com/sporkmonger/ecsevent"
)

func init() {
	zerologFieldInit()
}

func zerologFieldInit() {
	zerolog.LevelFieldName = ecsevent.FieldLogLevel
	zerolog.MessageFieldName = ecsevent.FieldMessage
	zerolog.ErrorFieldName = ecsevent.FieldErrorMessage
	zerolog.TimestampFieldName = ecsevent.FieldTimestamp
	zerolog.TimeFieldFormat = time.RFC3339Nano
	// Do not use caller for now.
}
