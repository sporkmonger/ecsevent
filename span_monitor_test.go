package ecsevent

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpanMonitor(t *testing.T) {
	mock := &mockEmitter{events: make([]map[string]interface{}, 0)}
	tcs := []struct {
		name                      string
		options                   []MonitorOption
		events                    []map[string]interface{}
		expectedEventsRecorded    int
		expectedSubeventsRecorded int
	}{
		{
			"no options set, no events sent",
			nil,
			[]map[string]interface{}{},
			0,
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
			0,
		},
		{
			"mock emitter, no events sent",
			[]MonitorOption{
				EmitToMock(mock),
				NestEvents(false),
			},
			[]map[string]interface{}{},
			1,
			0,
		},
		{
			"mock emitter, 1 event sent",
			[]MonitorOption{
				EmitToMock(mock),
				NestEvents(false),
			},
			[]map[string]interface{}{
				map[string]interface{}{
					FieldMessage: "test message",
				},
			},
			1,
			1,
		},
		{
			"mock emitter, 3 events sent",
			[]MonitorOption{
				EmitToMock(mock),
				NestEvents(false),
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
			1,
			3,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			// reset mock's events to empty
			mock.events = make([]map[string]interface{}, 0)
			rm := New(tc.options...)
			sm := NewSpanMonitorFromParent(rm)
			for _, event := range tc.events {
				sm.Record(event)
			}
			sm.Finish()
			// assert.Equal(rm, sm.Root())
			assert.Len(mock.events, tc.expectedEventsRecorded)
			if len(mock.events) >= 1 {
				event := mock.events[0]
				subevents := event[FieldEventSubevents]
				if subevents == nil {
					assert.Equal(0, tc.expectedSubeventsRecorded)
				} else {
					assert.Len(event[FieldEventSubevents], tc.expectedSubeventsRecorded)
				}
			}
		})
	}
}

func TestSpanMonitorConcurrency(t *testing.T) {
	mock := &mockEmitter{events: make([]map[string]interface{}, 0)}
	event := map[string]interface{}{
		FieldMessage: "test message",
	}

	t.Run("span monitor concurrency", func(t *testing.T) {
		assert := assert.New(t)
		rm := New(EmitToMock(mock))
		sm := NewSpanMonitorFromParent(rm)
		wg := sync.WaitGroup{}

		// This reliably triggers a failure if mutexes are not set up correctly
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 5000; j++ {
					sm.Record(event)
				}
				sm.Finish()
				// assert.Equal(rm, sm.Root())
			}()
		}
		wg.Wait()
		assert.Equal(10, len(mock.events))
	})
}

func TestSpanMonitorNesting(t *testing.T) {
	mock := &mockEmitter{events: make([]map[string]interface{}, 0)}
	event := map[string]interface{}{
		FieldMessage: "test message",
	}

	t.Run("span monitor nesting", func(t *testing.T) {
		assert := assert.New(t)
		rm := New(EmitToMock(mock))
		sm1 := NewSpanMonitorFromParent(rm)
		sm2 := NewSpanMonitorFromParent(sm1)
		sm2.Record(event)
		sm2.Finish()
		sm1.Finish()

		assert.Equal(rm, sm1.Root())
		assert.Equal(rm, sm2.Root())
		assert.Equal(sm1, sm2.Parent())
		assert.Equal(1, len(mock.events))
	})
}
