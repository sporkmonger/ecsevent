package ecsevent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendStackDriver(t *testing.T) {
	tcs := []struct {
		name           string
		input          map[string]interface{}
		expectedOutput map[string]interface{}
	}{
		{
			"simple message",
			map[string]interface{}{
				FieldTimestamp: "2019-10-28T06:15:07.226113003Z",
				FieldLogLevel:  "info",
				FieldMessage:   "something happened",
			},
			map[string]interface{}{
				FieldTimestamp: "2019-10-28T06:15:07.226113003Z",
				FieldLogLevel:  "info",
				FieldMessage:   "something happened",
				"timestamp":    "2019-10-28T06:15:07.226113003Z",
				"severity":     "INFO",
			},
		},
		{
			"http request",
			map[string]interface{}{
				FieldHTTPRequestMethod:      "GET",
				FieldHTTPVersion:            "1.1",
				FieldHTTPResponseBytes:      42,
				FieldHTTPResponseStatusCode: 200,
				FieldURLPath:                "/",
				FieldURLDomain:              "example.com",
				FieldEventDuration:          50000000,
			},
			map[string]interface{}{
				FieldHTTPRequestMethod:      "GET",
				FieldHTTPVersion:            "1.1",
				FieldHTTPResponseBytes:      42,
				FieldHTTPResponseStatusCode: 200,
				FieldURLPath:                "/",
				FieldURLDomain:              "example.com",
				FieldEventDuration:          50000000,
				"httpRequest.requestMethod": "GET",
				"httpRequest.protocol":      "HTTP/1.1",
				"httpRequest.status":        200,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			output := appendStackDriver(tc.input)
			assert.Equal(tc.expectedOutput, output)
		})
	}
}

func TestStackDriverSeverity(t *testing.T) {
	tcs := []struct {
		level            string
		expectedSeverity string
	}{
		{"t", "DEBUG"},
		{"trace", "DEBUG"},
		{"d", "DEBUG"},
		{"DEBUG", "DEBUG"},
		{"debug", "DEBUG"},
		{"i", "INFO"},
		{"info", "INFO"},
		{"n", "NOTICE"},
		{"notice", "NOTICE"},
		{"w", "WARNING"},
		{"warn", "WARNING"},
		{"e", "ERROR"},
		{"error", "ERROR"},
		{"c", "CRITICAL"},
		{"critical", "CRITICAL"},
		{"a", "ALERT"},
		{"alert", "ALERT"},
		{"f", "EMERGENCY"},
		{"fatal", "EMERGENCY"},
		{"b", "DEFAULT"},
		{"bogus", "DEFAULT"},
	}

	for _, tc := range tcs {
		t.Run(tc.level, func(t *testing.T) {
			assert := assert.New(t)
			output := stackDriverSeverity(tc.level)
			assert.Equal(tc.expectedSeverity, output)
		})
	}
}
