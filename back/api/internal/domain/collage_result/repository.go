package collage_result

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Create creates a new collage result
	Create(ctx context.Context, result *CollageResult) error

	// FindByID finds a collage result by ID
	FindByID(ctx context.Context, resultID uuid.UUID) (*CollageResult, error)

	// FindByGroupID finds all collage results by group ID
	FindByGroupID(ctx context.Context, groupID string, limit, offset int) ([]*CollageResult, error)

	// FindUnnotified finds all unnotified collage results
	FindUnnotified(ctx context.Context, limit int) ([]*CollageResult, error)

	// Update updates a collage result
	Update(ctx context.Context, result *CollageResult) error

	// Delete deletes a collage result
	Delete(ctx context.Context, resultID uuid.UUID) error
}
