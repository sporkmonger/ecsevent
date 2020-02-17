package humio

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/sporkmonger/ecsevent"
)

// Emitter wraps an HTTP client allowing ECS formatted events to be
// emitted to Humio.
type Emitter struct {
	// Server is the origin to ship logs to. No trailing slash, typically.
	// Defaults to 'https://cloud.humio.com'.
	Server string
	// IngestToken is the required API key to send events to Humio.
	IngestToken string
	// Tags are optional key-value pairs associated with events.
	// They must be low-cardinality.
	Tags     map[string]string
	client   *http.Client
	endpoint *url.URL
}

type ingestEvent struct {
	Timestamp  time.Time              `json:"timestamp"`
	Attributes map[string]interface{} `json:"attributes"`
}

type ingestRequest struct {
	Tags   map[string]string `json:"tags"`
	Events []ingestEvent     `json:"events"`
}

// Emit takes a map of ECS fields and values and emits the event into the
// Beeline.
func (e *Emitter) Emit(event map[string]interface{}) {
	// TODO: maybe a constructor?
	// Quick and dirty for now.
	if e.client == nil {
		e.client = &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 4 * time.Second,
				ResponseHeaderTimeout: 3 * time.Second,
			},
			Timeout: 2 * time.Minute,
		}
	}
	if e.Server == "" {
		e.Server = "https://cloud.humio.com"
	}
	if e.endpoint == nil {
		e.endpoint, _ = url.Parse(e.Server + "/api/v1/ingest/humio-structured")
	}
	if e.Tags == nil {
		e.Tags = make(map[string]string)
	}
	// Humio seems to recommend against nesting, but supports it. Nest for now,
	// maybe make this configurable.
	// unnested := ecsevent.Unnest(event)
	timestamp, ok := event[ecsevent.FieldTimestamp].(time.Time)
	if !ok {
		timestampString, ok := event[ecsevent.FieldTimestamp].(string)
		if ok {
			timestamp, _ = time.Parse(time.RFC3339Nano, timestampString)
		}
	}
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	ir := []ingestRequest{
		ingestRequest{
			Tags: e.Tags,
			Events: []ingestEvent{
				ingestEvent{
					Timestamp:  timestamp,
					Attributes: event,
				},
			},
		},
	}
	data, _ := json.Marshal(ir)
	b := bytes.NewBuffer(data)
	req, _ := http.NewRequest(http.MethodPost, e.endpoint.String(), b)
	// TODO: Emit _really_ needs to return an error
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.IngestToken)
	resp, err := e.client.Do(req)
	if err == nil && resp != nil {
		defer resp.Body.Close()
	}
}

var (
	// This is a compile-time check to make sure our types correctly
	// implement the interface:
	// https://medium.com/@matryer/c167afed3aae
	_ ecsevent.Emitter = &Emitter{}
)
