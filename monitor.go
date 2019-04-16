package ecsevent

import (
	"github.com/opentracing/opentracing-go"
)

var disabledMonitor *GlobalMonitor

func init() {
	m := Nop()
	disabledMonitor = m
}

type Monitor interface {
	Disabled() bool
	Fields() map[string]interface{}
	// StartSpan(name string, opts ...opentracing.StartSpanOption) (opentracing.Span, Monitor)
}

type GlobalMonitor struct {
	// Fields are the globally scoped fields applied to all events recorded by
	// the logger.
	fields   map[string]interface{}
	emitter  Emitter
	tracer   opentracing.Tracer
	disabled bool
}

var (
	// This is a compile-time check to make sure our types correctly
	// implement the interface:
	// https://medium.com/@matryer/c167afed3aae
	_ Monitor = &GlobalMonitor{}
)

// MonitorOptions configure a GlobalMonitor as it's being initialized.
type MonitorOptions func(*GlobalMonitor)

func New(opts ...MonitorOptions) Monitor {
	monitor := &GlobalMonitor{fields: make(map[string]interface{})}
	for _, opts := range opts {
		opts(monitor)
	}
	return monitor
}

func (gm *GlobalMonitor) Disabled() bool {
	return gm.disabled
}

func (gm *GlobalMonitor) Fields() map[string]interface{} {
	return gm.fields
}

// Nop returns a disabled logger for which all operation are no-op.
func Nop() *GlobalMonitor {
	return &GlobalMonitor{disabled: true}
}
