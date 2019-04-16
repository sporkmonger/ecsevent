package ecsevent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNest(t *testing.T) {
	tcs := []struct {
		name           string
		input          map[string]interface{}
		expectedOutput map[string]interface{}
	}{
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
				"http": map[string]interface{}{
					"request": map[string]interface{}{
						"method": "GET",
					},
					"response": map[string]interface{}{
						"bytes":       42,
						"status_code": 200,
					},
					"version": "1.1",
				},
				"url": map[string]interface{}{
					"domain": "example.com",
					"path":   "/",
				},
				"event": map[string]interface{}{
					"duration": 50000000,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			output := Nest(tc.input)
			assert.Equal(tc.expectedOutput, output)
		})
	}
}

func BenchmarkNest(b *testing.B) {
	input := map[string]interface{}{
		FieldHTTPRequestMethod:      "GET",
		FieldHTTPVersion:            "1.1",
		FieldHTTPResponseBytes:      42,
		FieldHTTPResponseStatusCode: 200,
		FieldURLPath:                "/",
		FieldURLDomain:              "example.com",
		FieldEventDuration:          50000000,
	}

	for n := 0; n < b.N; n++ {
		Nest(input)
	}
}
