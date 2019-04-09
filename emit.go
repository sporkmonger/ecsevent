package slog

// Emitter is a common interface for all SLog adapters.
type Emitter interface {
	// Emit takes a flat map of ECS fields and values, converts it to a nested
	// map, and emits the event on the underlying logger implementation.
	Emit(event map[string]interface{})
}
