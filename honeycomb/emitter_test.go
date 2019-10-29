package honeycomb

import (
	"testing"

	"github.com/sporkmonger/ecsevent"

	libhoney "github.com/honeycombio/libhoney-go"
	"github.com/stretchr/testify/assert"
)

func TestEmitter(t *testing.T) {
	tcs := []struct {
		name  string
		event map[string]interface{}
	}{
		{
			"empty event",
			map[string]interface{}{},
		},
		{
			"minimal event",
			map[string]interface{}{
				ecsevent.FieldMessage: "hello world",
			},
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
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			client, err := libhoney.NewClient(libhoney.ClientConfig{})
			assert.NoError(err)
			emitter := Emitter{Client: client}
			// won't actually send, no write key
			emitter.Emit(tc.event)
		})
	}
}
