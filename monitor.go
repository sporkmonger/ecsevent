package ecsevent

import (
	"sync"

	"github.com/opentracing/opentracing-go"
)

type Monitor interface {
	Fields() map[string]interface{}
	UpdateFields(map[string]interface{})
	Record(map[string]interface{})
	Root() *RootMonitor
}

type syncEmitter struct {
	emitter Emitter
	// mu gates events emitted since we don't expect emitters to be thread-safe
	mu sync.Mutex
}

// Emit takes a flat map of ECS fields and values, converts it to a nested
// map, and emits the event on the underlying logger implementation.
func (se syncEmitter) Emit(event map[string]interface{}) {
	se.emitter.Emit(event)
}

type RootMonitor struct {
	// Fields are the globally scoped fields applied to all events recorded by
	// the logger.
	fields      map[string]interface{}
	emitters    []*syncEmitter
	tracer      opentracing.Tracer
	nested      bool
	stackdriver bool
	// mu gates everything in this struct, including changes to the emitter
	// list but not events being emitted by the emitters.
	mu sync.Mutex
}

var (
	// This is a compile-time check to make sure our types correctly
	// implement the interface:
	// https://medium.com/@matryer/c167afed3aae
	_ Monitor = &RootMonitor{}
	_ Emitter = syncEmitter{}
)

// MonitorOption configure a RootMonitor as it's being initialized.
type MonitorOption func(*RootMonitor)

// NestEvents controls whether event fields should be nested or left
// in dot-notated format.
func NestEvents(nested bool) MonitorOption {
	return func(rm *RootMonitor) {
		rm.nested = nested
	}
}

// Stackdriver controls whether ECS events will automatically convert to
// the special fields expected by Stackdriver.
//
// The original ECS values will still be logged, but new Stackdriver values
// will also be created in the expected fields with any necessary transforms
// applied.
func Stackdriver(stackdriver bool) MonitorOption {
	return func(rm *RootMonitor) {
		rm.stackdriver = stackdriver
	}
}

// New creates a new RootMonitor with the given MonitorOption functions
// applied.
func New(opts ...MonitorOption) Monitor {
	return NewRootMonitor(opts...)
}

// NewRootMonitor creates a new RootMonitor with the given MonitorOption functions
// applied.
func NewRootMonitor(opts ...MonitorOption) *RootMonitor {
	monitor := &RootMonitor{
		fields:   make(map[string]interface{}),
		emitters: make([]*syncEmitter, 0),
		// avoid unneeded nil checks
		tracer: opentracing.NoopTracer{},
		nested: true,
	}
	for _, opts := range opts {
		opts(monitor)
	}
	return monitor
}

// Root returns itself.
func (rm *RootMonitor) Root() *RootMonitor {
	return rm
}

// AppendEmitter adds an emitter to the RootMonitor's emitter list.
//
// This function is intended to be used inside of a MonitorOption function
// and generally should not be used outside of initialization.
func (rm *RootMonitor) AppendEmitter(emitter Emitter) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.emitters = append(rm.emitters, &syncEmitter{emitter: emitter})
}

// SetTracer sets the tracer for the RootMonitor. Unlike emitters, there
// can be only one tracer.
//
// This function is intended to be used inside of a MonitorOption function
// and generally should not be used outside of initialization.
func (rm *RootMonitor) SetTracer(tracer opentracing.Tracer) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.tracer = tracer
}

// SetStackdriverLogging enables or disables translation of ECS events into
// the fields needed by Stackdriver.
func (rm *RootMonitor) SetStackdriverLogging(enabled bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.stackdriver = enabled
}

// Fields returns the fields currently set on the monitor.
func (rm *RootMonitor) Fields() map[string]interface{} {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	return rm.fields
}

// UpdateFields updates the RootMonitor's Field set.
func (rm *RootMonitor) UpdateFields(fields map[string]interface{}) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.fields == nil {
		rm.fields = make(map[string]interface{})
	}
	for k, v := range fields {
		rm.fields[k] = v
	}
}

// Record takes a series of fields and records an event.
func (rm *RootMonitor) Record(event map[string]interface{}) {
	for _, se := range rm.emitters {
		se.mu.Lock()
		// TODO: use fields
		if rm.stackdriver {
			event = appendStackdriver(event)
		}
		if rm.nested {
			event = Nest(event)
		}
		se.emitter.Emit(event)
		se.mu.Unlock()
	}
}
