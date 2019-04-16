package ecsevent

import (
	"github.com/opentracing/opentracing-go"
)

// SpanMonitor is a short-lived monitor with additional contextual fields.
// It's typically used in conjunction with a Context. It relies on a parent
// Monitor to emit.
type SpanMonitor struct {
	// SubeventsField is the field that all subevents will be recorded under.
	// If no SubeventsField is set, subevents will be recorded to emitters as
	// separate events. Has no effect on subevents emitted to an opentracing
	// span.
	SubeventsField string

	// subevents tracks all events recorded against this span monitor.
	subevents []map[string]interface{}

	// Fields are the locally scoped fields applied to all events and subevents
	// recorded by the SpanMonitor.
	fields map[string]interface{}

	// The opentracing span, if any, associated with this SpanMonitor.
	span   opentracing.Span
	parent Monitor
}

var (
	// This is a compile-time check to make sure our types correctly
	// implement the interface:
	// https://medium.com/@matryer/c167afed3aae
	_ Monitor = &SpanMonitor{}
)

func (sm *SpanMonitor) Fields() map[string]interface{} {
	return sm.fields
}

// Parent returns the parent logger.
func (sm SpanMonitor) Parent() Monitor {
	return sm.parent
}

// UpdateFields updates the SpanMonitor's field set.
func (sm SpanMonitor) UpdateFields(fields map[string]interface{}) {
	if sm.fields == nil {
		sm.fields = make(map[string]interface{})
	}
	for k, v := range fields {
		sm.fields[k] = v
	}
}

// Record takes a series of fields and records an event.
func (sm SpanMonitor) Record(event map[string]interface{}) {
	if sm.fields == nil {
		sm.fields = make(map[string]interface{})
	}
	merged := make(map[string]interface{})
	for k, v := range sm.fields {
		merged[k] = v
	}
	for k, v := range event {
		merged[k] = v
	}
	sm.subevents = append(sm.subevents, merged)
	// TODO: if configured to flush immediately, emit to parent, otherwise emit on Finish
}

func (sm SpanMonitor) Finish() {
	if sm.span != nil {
		records := make([]opentracing.LogRecord, 0)
		// TODO: generate log records for each subevent
		opts := opentracing.FinishOptions{LogRecords: records}
		sm.span.FinishWithOptions(opts)
	}
}
