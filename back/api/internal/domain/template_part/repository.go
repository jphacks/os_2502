package template_part

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, part *TemplatePart) error
	FindByID(ctx context.Context, partID uuid.UUID) (*TemplatePart, error)
	FindByTemplateID(ctx context.Context, templateID uuid.UUID) ([]*TemplatePart, error)
	Delete(ctx context.Context, partID uuid.UUID) error
}
