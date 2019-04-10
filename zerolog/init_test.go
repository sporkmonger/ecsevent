package zerolog

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/sporkmonger/ecslog"
)

func TestInit(t *testing.T) {
	assert := assert.New(t)

	zerologFieldInit()
	assert.Equal(ecslog.FieldLogLevel, zerolog.LevelFieldName)
	assert.Equal(ecslog.FieldMessage, zerolog.MessageFieldName)
	assert.Equal(ecslog.FieldErrorMessage, zerolog.ErrorFieldName)
	assert.Equal(ecslog.FieldTimestamp, zerolog.TimestampFieldName)
}
