package slog

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTypeCheck(t *testing.T) {
	tcs := []struct {
		fieldName     string
		value         interface{}
		expectedError error
	}{
		{
			FieldAgentEphemeralID,
			1234,
			errors.New("unexpected value for field 'agent.ephemeral_id', refer to ECS specification"),
		},
		{
			FieldTimestamp,
			1234,
			errors.New("unexpected value for field '@timestamp', refer to ECS specification"),
		},
		{
			FieldTimestamp,
			"2009-11-10T23:00:00Z", // wants a time.Time, not a string containing the time
			errors.New("unexpected value for field '@timestamp', refer to ECS specification"),
		},
		{
			FieldTimestamp,
			time.Now(),
			nil,
		},
		{
			FieldLabels,
			"a b c",
			errors.New("unexpected value for field 'labels', refer to ECS specification"),
		},
		{
			FieldLabels,
			[]string{"a", "b", "c"},
			errors.New("unexpected value for field 'labels', refer to ECS specification"),
		},
		{
			FieldLabels,
			map[string]string{"a": "a", "b": "b", "c": "c"},
			nil,
		},
		{
			FieldTags,
			"a b c",
			errors.New("unexpected value for field 'tags', refer to ECS specification"),
		},
		{
			FieldTags,
			[]string{"a", "b", "c"},
			nil,
		},
		{
			FieldTags,
			map[string]string{"a": "a", "b": "b", "c": "c"},
			errors.New("unexpected value for field 'tags', refer to ECS specification"),
		},
		{
			FieldHTTPResponseStatusCode,
			http.StatusText(http.StatusNotFound),
			errors.New("unexpected value for field 'http.response.status_code', refer to ECS specification"),
		},
		{
			FieldHTTPResponseStatusCode,
			http.StatusNotFound,
			nil,
		},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%s with %v", tc.fieldName, tc.value), func(t *testing.T) {
			assert := assert.New(t)
			err := typeCheck(tc.fieldName, tc.value)
			if tc.expectedError != nil {
				assert.EqualError(err, tc.expectedError.Error())
			} else {
				assert.NoError(err)
			}
		})
	}
}
