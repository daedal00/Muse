// backend/graph/context.go
package graph

import (
	"context"
)

// unexported type to avoid collisions with other packages
type ctxKey string

// UserIDKey stores userID, resolvers can access it directly now
const UserIDKey ctxKey = "userID"

// helper function to pull userID from context
func ForContext(ctx context.Context) (string, bool) {
	raw := ctx.Value(UserIDKey)
	if raw == nil {
		return "", false
	}
	id, ok := raw.(string)
	return id, ok
}