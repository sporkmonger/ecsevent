package honeycomb

import (
	libhoney "github.com/honeycombio/libhoney-go"

	"github.com/sporkmonger/ecsevent"
)

// Emitter wraps a Beeline allowing ECS formatted events to be emitted.
type Emitter struct {
	Client *libhoney.Client
}

// Emit takes a map of ECS fields and values and emits the event into the
// Beeline.
func (e *Emitter) Emit(event map[string]interface{}) {
	unnested := ecsevent.Unnest(event)
	he := e.Client.NewEvent()
	he.Add(unnested)
	he.Send()
}

var (
	// This is a compile-time check to make sure our types correctly
	// implement the interface:
	// https://medium.com/@matryer/c167afed3aae
	_ ecsevent.Emitter = &Emitter{}
)
