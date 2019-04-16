package ecsevent

import (
	"context"
)

// ctxKey is an unexported key for retrieving a logger from a context
type ctxKey struct{}

func (gm *GlobalMonitor) WithContext(ctx context.Context) context.Context {
	return NewContext(ctx, gm)
}

func (sm *SpanMonitor) WithContext(ctx context.Context) context.Context {
	return NewContext(ctx, sm)
}

func NewContext(ctx context.Context, m Monitor) context.Context {
	if mp, ok := ctx.Value(ctxKey{}).(Monitor); ok {
		if mp == m {
			// Do not store same monitor.
			return ctx
		}
	}
	return context.WithValue(ctx, ctxKey{}, m)
}

func FromContext(ctx context.Context) Monitor {
	if l, ok := ctx.Value(ctxKey{}).(Monitor); ok {
		return l
	}
	return disabledMonitor
}
