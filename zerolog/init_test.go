package zerolog

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/sporkmonger/ecsevent"
)

func TestInit(t *testing.T) {
	assert := assert.New(t)

	zerologFieldInit()
	assert.Equal(ecsevent.FieldLogLevel, zerolog.LevelFieldName)
	assert.Equal(ecsevent.FieldMessage, zerolog.MessageFieldName)
	assert.Equal(ecsevent.FieldErrorMessage, zerolog.ErrorFieldName)
	assert.Equal(ecsevent.FieldTimestamp, zerolog.TimestampFieldName)
}
