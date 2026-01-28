package graph

import "fmt"

func resolveLimit(first *int32, fallback int) int {
	if first != nil && *first > 0 {
		return int(*first)
	}
	return fallback
}

func (r *Resolver) resolveOffset(after *string) (int, error) {
	if after == nil || *after == "" {
		return 0, nil
	}

	cursor, err := r.paginationHelper.DecodeCursor(*after)
	if err != nil {
		return 0, fmt.Errorf("invalid cursor: %w", err)
	}

	if cursor.Position < 0 {
		return 0, nil
	}

	return cursor.Position + 1, nil
}
