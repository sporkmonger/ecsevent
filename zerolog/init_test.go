package zerolog

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/sporkmonger/slog"
)

func TestInit(t *testing.T) {
	assert := assert.New(t)

	zerologFieldInit()
	assert.Equal(slog.FieldLogLevel, zerolog.LevelFieldName)
	assert.Equal(slog.FieldMessage, zerolog.MessageFieldName)
	assert.Equal(slog.FieldErrorMessage, zerolog.ErrorFieldName)
	assert.Equal(slog.FieldTimestamp, zerolog.TimestampFieldName)
}
