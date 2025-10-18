package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/template_part"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db/models"
)

type TemplatePartRepositorySQLBoiler struct {
	db *sql.DB
}

func NewTemplatePartRepositorySQLBoiler(db *sql.DB) template_part.Repository {
	return &TemplatePartRepositorySQLBoiler{db: db}
}

func toTemplatePartModel(tp *template_part.TemplatePart) *models.TemplatePart {
	model := &models.TemplatePart{
		PartID:      tp.PartID().String(),
		TemplateID:  tp.TemplateID().String(),
		PartNumber:  tp.PartNumber(),
		PositionX:   tp.PositionX(),
		PositionY:   tp.PositionY(),
		Width:       tp.Width(),
		Height:      tp.Height(),
		CreatedAt:   tp.CreatedAt(),
		UpdatedAt:   tp.UpdatedAt(),
	}

	if partName := tp.PartName(); partName != nil {
		model.PartName = null.StringFrom(*partName)
	}

	if description := tp.Description(); description != nil {
		model.Description = null.StringFrom(*description)
	}

	return model
}

func toTemplatePartEntity(m *models.TemplatePart) (*template_part.TemplatePart, error) {
	partID, err := uuid.Parse(m.PartID)
	if err != nil {
		return nil, err
	}
	templateID, err := uuid.Parse(m.TemplateID)
	if err != nil {
		return nil, err
	}

	var partName *string
	if m.PartName.Valid {
		partName = &m.PartName.String
	}

	var description *string
	if m.Description.Valid {
		description = &m.Description.String
	}

	return template_part.Reconstruct(
		partID,
		templateID,
		m.PartNumber,
		m.PositionX,
		m.PositionY,
		m.Width,
		m.Height,
		partName,
		description,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

func (r *TemplatePartRepositorySQLBoiler) Save(ctx context.Context, part *template_part.TemplatePart) error {
	model := toTemplatePartModel(part)
	return model.Upsert(ctx, r.db, true, []string{"part_id"}, boil.Infer(), boil.Infer())
}

func (r *TemplatePartRepositorySQLBoiler) FindByID(ctx context.Context, partID uuid.UUID) (*template_part.TemplatePart, error) {
	model, err := models.FindTemplatePart(ctx, r.db, partID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, template_part.ErrNotFound
		}
		return nil, err
	}
	return toTemplatePartEntity(model)
}

func (r *TemplatePartRepositorySQLBoiler) FindByTemplateID(ctx context.Context, templateID uuid.UUID) ([]*template_part.TemplatePart, error) {
	modelSlice, err := models.TemplateParts(
		qm.Where("template_id = ?", templateID.String()),
		qm.OrderBy("part_number ASC"),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	entities := make([]*template_part.TemplatePart, len(modelSlice))
	for i, model := range modelSlice {
		entity, err := toTemplatePartEntity(model)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}

func (r *TemplatePartRepositorySQLBoiler) Delete(ctx context.Context, partID uuid.UUID) error {
	model, err := models.FindTemplatePart(ctx, r.db, partID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return template_part.ErrNotFound
		}
		return err
	}
	_, err = model.Delete(ctx, r.db)
	return err
}
