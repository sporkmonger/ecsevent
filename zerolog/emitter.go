package zerolog

import (
	"github.com/rs/zerolog"
	"github.com/sporkmonger/ecslog"
)

// Emitter wraps a zerolog Logger allowing ECS formatted events to be emitted.
type Emitter struct {
	Logger zerolog.Logger
}

// Emit takes a flat map of ECS fields and values, converts it to a nested
// map, and emits the event into the zerolog logger.
func (e *Emitter) Emit(event map[string]interface{}) {
	nested := ecslog.Nest(event)
	e.Logger.Log().Fields(nested).Msg("")
}

var (
	// This is a compile-time check to make sure our types correctly
	// implement the interface:
	// https://medium.com/@matryer/c167afed3aae
	_ ecslog.Emitter = &Emitter{}
)
