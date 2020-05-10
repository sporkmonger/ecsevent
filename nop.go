package ecsevent

var disabledMonitor *NopMonitor

func init() {
	m := Nop()
	disabledMonitor = m
}

// NopMonitor is a disabled monitor for which all operation are no-op.
type NopMonitor struct {
}

var (
	// This is a compile-time check to make sure our types correctly
	// implement the interface:
	// https://medium.com/@matryer/c167afed3aae
	_ Monitor = &NopMonitor{}
)

// Fields returns a throw-away, always empty map of the fields.
func (nm *NopMonitor) Fields() map[string]interface{} {
	return make(map[string]interface{})
}

// Record does nothing.
func (nm *NopMonitor) Record(event map[string]interface{}) {}

// UpdateFields does nothing.
func (nm *NopMonitor) UpdateFields(event map[string]interface{}) {}

// Nop returns a disabled monitor for which all operation are no-op.
func Nop() *NopMonitor {
	return &NopMonitor{}
}

// Root returns nil.
func (nm *NopMonitor) Root() *RootMonitor {
	return nil
}
