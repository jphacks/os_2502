package collage_template

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Create creates a new collage template
	Create(ctx context.Context, template *CollageTemplate) error

	// FindByID finds a collage template by ID
	FindByID(ctx context.Context, templateID uuid.UUID) (*CollageTemplate, error)

	// FindByName finds a collage template by name
	FindByName(ctx context.Context, name string) (*CollageTemplate, error)

	// List lists all collage templates
	List(ctx context.Context, limit, offset int) ([]*CollageTemplate, error)

	// Update updates a collage template
	Update(ctx context.Context, template *CollageTemplate) error

	// Delete deletes a collage template
	Delete(ctx context.Context, templateID uuid.UUID) error
}
