package ecsevent

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockEmitter struct {
	events []map[string]interface{}
}

func (me *mockEmitter) Emit(fields map[string]interface{}) {
	me.events = append(me.events, fields)
}

func (me *mockEmitter) Events() []map[string]interface{} {
	return me.events
}

func EmitToMock(mock *mockEmitter) MonitorOption {
	return func(gm *GlobalMonitor) {
		gm.AppendEmitter(mock)
	}
}

func TestNewGlobalMonitor(t *testing.T) {
	mock := &mockEmitter{events: make([]map[string]interface{}, 0)}
	tcs := []struct {
		name                   string
		options                []MonitorOption
		events                 []map[string]interface{}
		expectedEventsRecorded int
	}{
		{
			"no options set, no events sent",
			nil,
			[]map[string]interface{}{},
			0,
		},
		{
			"no options set, 1 event sent",
			nil,
			[]map[string]interface{}{
				map[string]interface{}{
					FieldMessage: "test message",
				},
			},
			0,
		},
		{
			"mock emitter, no events sent",
			[]MonitorOption{
				EmitToMock(mock),
			},
			[]map[string]interface{}{},
			0,
		},
		{
			"mock emitter, 1 event sent",
			[]MonitorOption{
				EmitToMock(mock),
			},
			[]map[string]interface{}{
				map[string]interface{}{
					FieldMessage: "test message",
				},
			},
			1,
		},
		{
			"mock emitter, 3 events sent",
			[]MonitorOption{
				EmitToMock(mock),
			},
			[]map[string]interface{}{
				map[string]interface{}{
					FieldMessage: "test message",
				},
				map[string]interface{}{
					FieldMessage: "another test message",
				},
				map[string]interface{}{
					FieldMessage: "and again",
				},
			},
			3,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			// reset mock's events to empty
			mock.events = make([]map[string]interface{}, 0)
			gm := New(tc.options...)
			for _, event := range tc.events {
				gm.Record(event)
			}
			assert.Len(mock.events, tc.expectedEventsRecorded)
		})
	}
}

func TestGlobalMonitorConcurrency(t *testing.T) {
	mock := &mockEmitter{events: make([]map[string]interface{}, 0)}
	event := map[string]interface{}{
		FieldMessage: "test message",
	}

	t.Run("global monitor concurrency", func(t *testing.T) {
		assert := assert.New(t)
		gm := New(EmitToMock(mock))
		wg := sync.WaitGroup{}
		// This reliably triggers a failure if mutexes are not set up correctly
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 5000; j++ {
					gm.Record(event)
				}
			}()
		}
		wg.Wait()
		assert.Equal(50000, len(mock.events))
	})
}
