package ecsevent

import (
	"context"
)

// ctxKey is an unexported key for retrieving a logger from a context
type ctxKey struct{}

func (rm *RootMonitor) WithContext(ctx context.Context) context.Context {
	return NewContext(ctx, rm)
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

func MonitorFromContext(ctx context.Context) Monitor {
	if m, ok := ctx.Value(ctxKey{}).(Monitor); ok {
		return m
	}
	return disabledMonitor
}
