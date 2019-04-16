package zerolog

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sporkmonger/ecsevent"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestEmitter(t *testing.T) {
	tcs := []struct {
		name           string
		event          map[string]interface{}
		expectedOutput string
	}{
		{
			"empty event",
			map[string]interface{}{},
			`{}`,
		},
		{
			"minimal event",
			map[string]interface{}{
				ecsevent.FieldMessage: "hello world",
			},
			`{"message":"hello world"}`,
		},
		{
			"http event",
			map[string]interface{}{
				ecsevent.FieldHTTPRequestMethod:      "GET",
				ecsevent.FieldHTTPRequestReferrer:    "http://example.com/",
				ecsevent.FieldHTTPResponseStatusCode: 200,
				ecsevent.FieldHTTPVersion:            "1.1",
				ecsevent.FieldSourceIP:               "127.0.0.1",
				ecsevent.FieldURLDomain:              "example.com",
				ecsevent.FieldURLFull:                "http://example.com/hello",
			},
			`{"http":{"request":{"method":"GET","referrer":"http://example.com/"},"response":{"status_code":200},"version":"1.1"},"source":{"ip":"127.0.0.1"},"url":{"domain":"example.com","full":"http://example.com/hello"}}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			buffer := &bytes.Buffer{}
			logger := zerolog.New(buffer)
			emitter := Emitter{Logger: logger}
			emitter.Emit(tc.event)
			output := buffer.String()
			assert.Equal(tc.expectedOutput, strings.TrimSpace(output))
		})
	}
}
