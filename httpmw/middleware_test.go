package httpmw

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sporkmonger/ecsevent"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerParent(t *testing.T) {
	assert := assert.New(t)
	monitor := ecsevent.New()
	mh := NewHandler(monitor)
	h := mh(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sm := FromRequest(r)
		assert.NotNil(sm)
		if sm != nil {
			assert.Equal(monitor, sm.Parent())
		}
	}))
	h.ServeHTTP(nil, &http.Request{})
}

type mockEmitter struct {
	events []map[string]interface{}
}

func (me *mockEmitter) Emit(fields map[string]interface{}) {
	me.events = append(me.events, fields)
}

func (me *mockEmitter) Events() []map[string]interface{} {
	return me.events
}

func EmitToMock(mock *mockEmitter) ecsevent.MonitorOption {
	return func(gm *ecsevent.GlobalMonitor) {
		gm.AppendEmitter(mock)
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"status": "ok"}`)
}

func TestHealthCheckHandlerUnnested(t *testing.T) {
	assert := assert.New(t)

	mock := &mockEmitter{events: make([]map[string]interface{}, 0)}
	monitor := ecsevent.New(EmitToMock(mock), ecsevent.NestEvents(false))
	mh := NewHandler(monitor)

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1:54321"
	req.Header.Set("User-Agent", "go-test/1.0")

	rr := httptest.NewRecorder()
	handler := mh(http.HandlerFunc(HealthCheckHandler))

	handler.ServeHTTP(rr, req)

	assert.Equal(http.StatusOK, rr.Code)
	assert.Equal(`{"status": "ok"}`, rr.Body.String())
	assert.Len(mock.events, 1)
	if len(mock.events) == 1 {
		expectedEvent := map[string]interface{}{
			"ecs.version":               "1.0.1",
			"client.ip":                 "127.0.0.1",
			"client.port":               54321,
			"http.request.body.bytes":   int64(0),
			"http.request.method":       "GET",
			"http.response.body.bytes":  int64(16),
			"http.response.status_code": 200,
			"http.version":              "1.1",
			"url.full":                  "/health-check",
			"url.original":              "/health-check",
			"url.path":                  "/health-check",
			"user_agent.original":       "go-test/1.0",
		}
		assert.Equal(expectedEvent, mock.events[0])
	}
}

func TestHealthCheckHandlerNested(t *testing.T) {
	assert := assert.New(t)

	mock := &mockEmitter{events: make([]map[string]interface{}, 0)}
	monitor := ecsevent.New(EmitToMock(mock), ecsevent.NestEvents(true))
	mh := NewHandler(monitor)

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1:54321"
	req.Header.Set("User-Agent", "go-test/1.0")

	rr := httptest.NewRecorder()
	handler := mh(http.HandlerFunc(HealthCheckHandler))

	handler.ServeHTTP(rr, req)

	assert.Equal(http.StatusOK, rr.Code)
	assert.Equal(`{"status": "ok"}`, rr.Body.String())
	assert.Len(mock.events, 1)
	if len(mock.events) == 1 {
		expectedEvent := map[string]interface{}{
			"ecs": map[string]interface{}{
				"version": "1.0.1",
			},
			"client": map[string]interface{}{
				"ip":   "127.0.0.1",
				"port": 54321,
			},
			"http": map[string]interface{}{
				"request": map[string]interface{}{
					"body": map[string]interface{}{
						"bytes": int64(0),
					},
					"method": "GET",
				},
				"response": map[string]interface{}{
					"body": map[string]interface{}{
						"bytes": int64(16),
					},
					"status_code": 200,
				},
				"version": "1.1",
			},
			"url": map[string]interface{}{
				"full":     "/health-check",
				"original": "/health-check",
				"path":     "/health-check",
			},
			"user_agent": map[string]interface{}{
				"original": "go-test/1.0",
			},
		}
		assert.Equal(expectedEvent, mock.events[0])
	}
}
