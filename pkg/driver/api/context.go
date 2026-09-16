package api

import (
	"context"
)

type providerIDKey struct{}

// WithProviderID attaches the API driver / provider id to the context for auth lookup.
func WithProviderID(ctx context.Context, providerID string) context.Context {
	if providerID == "" {
		return ctx
	}
	return context.WithValue(ctx, providerIDKey{}, providerID)
}

func providerIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v, _ := ctx.Value(providerIDKey{}).(string)
	return v
}
