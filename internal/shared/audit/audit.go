package audit

import "context"

type ctxKey struct{}

// WithActor returns a context carrying the id of the acting user, used by the
// base model hooks to fill created_by / updated_by.
func WithActor(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, userID)
}

func Actor(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(ctxKey{}).(string); ok {
		return v
	}
	return ""
}
