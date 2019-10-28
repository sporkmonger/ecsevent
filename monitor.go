package ecsevent

import (
	"sync"

	"github.com/opentracing/opentracing-go"
)

type Monitor interface {
	Fields() map[string]interface{}
	UpdateFields(map[string]interface{})
	Record(map[string]interface{})
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

type GlobalMonitor struct {
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
	_ Monitor = &GlobalMonitor{}
	_ Emitter = syncEmitter{}
)

// MonitorOption configure a GlobalMonitor as it's being initialized.
type MonitorOption func(*GlobalMonitor)

// NestEvents controls whether event fields should be nested or left
// in dot-notated format.
func NestEvents(nested bool) MonitorOption {
	return func(gm *GlobalMonitor) {
		gm.nested = nested
	}
}

// Stackdriver controls whether ECS events will automatically convert to
// the special fields expected by Stackdriver.
//
// The original ECS values will still be logged, but new Stackdriver values
// will also be created in the expected fields with any necessary transforms
// applied.
func Stackdriver(stackdriver bool) MonitorOption {
	return func(gm *GlobalMonitor) {
		gm.stackdriver = stackdriver
	}
}

// New creates a new GlobalMonitor with the given MonitorOption functions
// applied.
func New(opts ...MonitorOption) Monitor {
	monitor := &GlobalMonitor{
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

// AppendEmitter adds an emitter to the GlobalMonitor's emitter list.
//
// This function is intended to be used inside of a MonitorOption function
// and generally should not be used outside of initialization.
func (gm *GlobalMonitor) AppendEmitter(emitter Emitter) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	gm.emitters = append(gm.emitters, &syncEmitter{emitter: emitter})
}

// SetTracer sets the tracer for the GlobalMonitor. Unlike emitters, there
// can be only one tracer.
//
// This function is intended to be used inside of a MonitorOption function
// and generally should not be used outside of initialization.
func (gm *GlobalMonitor) SetTracer(tracer opentracing.Tracer) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	gm.tracer = tracer
}

func (gm *GlobalMonitor) SetStackdriverLogging(enabled bool) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	gm.stackdriver = enabled
}

func (gm *GlobalMonitor) Fields() map[string]interface{} {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	return gm.fields
}

// UpdateFields updates the GlobalMonitor's Field set.
func (gm *GlobalMonitor) UpdateFields(fields map[string]interface{}) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	if gm.fields == nil {
		gm.fields = make(map[string]interface{})
	}
	for k, v := range fields {
		gm.fields[k] = v
	}
}

func (gm *GlobalMonitor) Record(event map[string]interface{}) {
	for _, se := range gm.emitters {
		se.mu.Lock()
		// TODO: use fields
		if gm.stackdriver {
			event = appendStackdriver(event)
		}
		if gm.nested {
			event = Nest(event)
		}
		se.emitter.Emit(event)
		se.mu.Unlock()
	}
}
