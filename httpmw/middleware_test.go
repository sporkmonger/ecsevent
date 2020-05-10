package httpmw

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/sporkmonger/ecsevent"

	"github.com/stretchr/testify/assert"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func TestNewHandlerParent(t *testing.T) {
	assert := assert.New(t)
	monitor := ecsevent.New()
	mh := NewHandler(monitor.(*ecsevent.GlobalMonitor))
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

	testStart := time.Now()

	mock := &mockEmitter{events: make([]map[string]interface{}, 0)}
	monitor := ecsevent.New(EmitToMock(mock), ecsevent.NestEvents(false))
	mh := NewHandler(monitor.(*ecsevent.GlobalMonitor))

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1:54321"
	req.Header.Set("User-Agent", "go-test/1.0")

	rr := httptest.NewRecorder()
	handler := mh(http.HandlerFunc(HealthCheckHandler))

	handler.ServeHTTP(rr, req)

	testEnd := time.Now()

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
		assert.Greater(mock.events[0]["event.duration"], int64(0))
		assert.Less(mock.events[0]["event.duration"], int64(1000000))
		delete(mock.events[0], "event.duration") // can't predict this value
		if eventStart, ok := mock.events[0]["event.start"].(time.Time); ok {
			assert.True(eventStart.After(testStart))
			assert.True(eventStart.Before(testEnd))
			if eventEnd, ok := mock.events[0]["event.end"].(time.Time); ok {
				assert.True(eventEnd.After(testStart))
				assert.True(eventEnd.Before(testEnd))
				assert.True(eventEnd.After(eventStart))
			} else {
				assert.Fail("event.end was not time.Time")
			}
		} else {
			assert.Fail("event.start was not time.Time")
		}
		delete(mock.events[0], "event.start") // can't predict this value
		delete(mock.events[0], "event.end")   // can't predict this value
		assert.Equal(expectedEvent, mock.events[0])
	}
}

func TestHealthCheckHandlerNested(t *testing.T) {
	assert := assert.New(t)

	testStart := time.Now()

	mock := &mockEmitter{events: make([]map[string]interface{}, 0)}
	monitor := ecsevent.New(EmitToMock(mock), ecsevent.NestEvents(true))
	mh := NewHandler(monitor.(*ecsevent.GlobalMonitor))

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1:54321"
	req.Header.Set("User-Agent", "go-test/1.0")

	rr := httptest.NewRecorder()
	handler := mh(http.HandlerFunc(HealthCheckHandler))

	handler.ServeHTTP(rr, req)

	testEnd := time.Now()

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
		eventMap := mock.events[0]["event"].(map[string]interface{})
		assert.Greater(eventMap["duration"], int64(0))
		assert.Less(eventMap["duration"], int64(1000000))
		if eventStart, ok := eventMap["start"].(time.Time); ok {
			assert.True(eventStart.After(testStart))
			assert.True(eventStart.Before(testEnd))
			if eventEnd, ok := eventMap["end"].(time.Time); ok {
				assert.True(eventEnd.After(testStart))
				assert.True(eventEnd.Before(testEnd))
				assert.True(eventEnd.After(eventStart))
			} else {
				assert.Fail("event.end was not time.Time")
			}
		} else {
			assert.Fail("event.start was not time.Time")
		}
		delete(mock.events[0], "event") // can't predict its values
		assert.Equal(expectedEvent, mock.events[0])
	}
}

func TestOpenTracing(t *testing.T) {
	assert := assert.New(t)

	tracer := mocktracer.New()
	opentracing.SetGlobalTracer(tracer)

	testStart := time.Now()
	time.Sleep(1 * time.Nanosecond)

	mock := &mockEmitter{events: make([]map[string]interface{}, 0)}
	monitor := ecsevent.New(EmitToMock(mock), ecsevent.NestEvents(true))
	mh := NewHandler(monitor.(*ecsevent.GlobalMonitor))

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "www.example.com"
	req.RemoteAddr = "127.0.0.1:54321"
	req.Header.Set("User-Agent", "go-test/1.0")

	parentSpan := opentracing.StartSpan("parent_span")
	opentracing.GlobalTracer().Inject(
		parentSpan.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	rr := httptest.NewRecorder()
	handler := mh(http.HandlerFunc(HealthCheckHandler))

	handler.ServeHTTP(rr, req)

	time.Sleep(1 * time.Nanosecond)
	testEnd := time.Now()

	assert.Equal(http.StatusOK, rr.Code)
	assert.Equal(`{"status": "ok"}`, rr.Body.String())

	parentSpan.Finish()

	finishedSpans := tracer.FinishedSpans()
	assert.Len(finishedSpans, 2, "both parent and child spans should be finished")
	if len(finishedSpans) != 1 {
		return
	}
	mockedSpan := finishedSpans[0]
	mockedParentSpan := finishedSpans[1]

	assert.True(mockedSpan.StartTime.After(testStart))
	assert.True(mockedSpan.StartTime.Before(testEnd))
	assert.True(mockedSpan.FinishTime.After(testStart))
	assert.True(mockedSpan.FinishTime.Before(testEnd))

	assert.NotEqual(0, mockedSpan.ParentID)
	assert.Equal(mockedParentSpan.SpanContext.SpanID, mockedSpan.ParentID)
	assert.Equal("GET www.example.com", mockedSpan.OperationName)
}

func TestOpenTracingWithJaeger(t *testing.T) {
	assert := assert.New(t)

	cfg := jaegercfg.Configuration{
		ServiceName: "test_service",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	reporter := jaeger.NewInMemoryReporter()
	tracer, closer, err := cfg.NewTracer(jaegercfg.Reporter(reporter))
	assert.NoError(err)
	if err != nil {
		return
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	mock := &mockEmitter{events: make([]map[string]interface{}, 0)}
	monitor := ecsevent.New(EmitToMock(mock), ecsevent.NestEvents(true))
	mh := NewHandler(monitor.(*ecsevent.GlobalMonitor))

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1:54321"
	req.Header.Set("User-Agent", "go-test/1.0")

	span := opentracing.StartSpan("parent_span")
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	rr := httptest.NewRecorder()
	handler := mh(http.HandlerFunc(HealthCheckHandler))

	handler.ServeHTTP(rr, req)

	assert.Equal(http.StatusOK, rr.Code)
	assert.Equal(`{"status": "ok"}`, rr.Body.String())

	assert.Equal(reporter.SpansSubmitted(), 1, "there should be only one submitted span")
}
