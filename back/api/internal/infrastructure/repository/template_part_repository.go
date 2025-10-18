package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/null/v8"
	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/template_part"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/models"
)

type TemplatePartRepository struct {
	db *sql.DB
}

func NewTemplatePartRepository(db *sql.DB) template_part.Repository {
	return &TemplatePartRepository{db: db}
}

// Model to Entity conversion
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

// Entity to Model conversion
func toTemplatePartModel(tp *template_part.TemplatePart) *models.TemplatePart {
	model := &models.TemplatePart{
		PartID:     tp.PartID().String(),
		TemplateID: tp.TemplateID().String(),
		PartNumber: tp.PartNumber(),
		PositionX:  tp.PositionX(),
		PositionY:  tp.PositionY(),
		Width:      tp.Width(),
		Height:     tp.Height(),
		CreatedAt:  tp.CreatedAt(),
		UpdatedAt:  tp.UpdatedAt(),
	}

	if partName := tp.PartName(); partName != nil {
		model.PartName = null.String{String: *partName, Valid: true}
	}

	if description := tp.Description(); description != nil {
		model.Description = null.String{String: *description, Valid: true}
	}

	return model
}

func (r *TemplatePartRepository) Create(ctx context.Context, tp *template_part.TemplatePart) error {
	model := toTemplatePartModel(tp)
	return model.Insert(ctx, r.db, boil.Infer())
}

func (r *TemplatePartRepository) FindByID(ctx context.Context, partID uuid.UUID) (*template_part.TemplatePart, error) {
	model, err := models.FindTemplatePart(ctx, r.db, partID.String())
	if err == sql.ErrNoRows {
		return nil, template_part.ErrTemplatePartNotFound
	}
	if err != nil {
		return nil, err
	}
	return toTemplatePartEntity(model)
}

func (r *TemplatePartRepository) FindByTemplateID(ctx context.Context, templateID uuid.UUID) ([]*template_part.TemplatePart, error) {
	modelSlice, err := models.TemplateParts(
		qm.Where("template_id = ?", templateID.String()),
		qm.OrderBy("part_number"),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	parts := make([]*template_part.TemplatePart, len(modelSlice))
	for i, model := range modelSlice {
		tp, err := toTemplatePartEntity(model)
		if err != nil {
			return nil, err
		}
		parts[i] = tp
	}
	return parts, nil
}

func (r *TemplatePartRepository) FindByTemplateIDAndPartNumber(ctx context.Context, templateID uuid.UUID, partNumber int) (*template_part.TemplatePart, error) {
	model, err := models.TemplateParts(
		qm.Where("template_id = ? AND part_number = ?", templateID.String(), partNumber),
	).One(ctx, r.db)
	if err == sql.ErrNoRows {
		return nil, template_part.ErrTemplatePartNotFound
	}
	if err != nil {
		return nil, err
	}
	return toTemplatePartEntity(model)
}

func (r *TemplatePartRepository) Update(ctx context.Context, tp *template_part.TemplatePart) error {
	model, err := models.FindTemplatePart(ctx, r.db, tp.PartID().String())
	if err == sql.ErrNoRows {
		return template_part.ErrTemplatePartNotFound
	}
	if err != nil {
		return err
	}

	model.PartNumber = tp.PartNumber()
	model.PositionX = tp.PositionX()
	model.PositionY = tp.PositionY()
	model.Width = tp.Width()
	model.Height = tp.Height()
	model.UpdatedAt = tp.UpdatedAt()

	if partName := tp.PartName(); partName != nil {
		model.PartName = null.String{String: *partName, Valid: true}
	} else {
		model.PartName = null.String{Valid: false}
	}

	if description := tp.Description(); description != nil {
		model.Description = null.String{String: *description, Valid: true}
	} else {
		model.Description = null.String{Valid: false}
	}

	_, err = model.Update(ctx, r.db, boil.Whitelist(
		models.TemplatePartColumns.PartNumber,
		models.TemplatePartColumns.PartName,
		models.TemplatePartColumns.PositionX,
		models.TemplatePartColumns.PositionY,
		models.TemplatePartColumns.Width,
		models.TemplatePartColumns.Height,
		models.TemplatePartColumns.Description,
		models.TemplatePartColumns.UpdatedAt,
	))
	return err
}

func (r *TemplatePartRepository) Delete(ctx context.Context, partID uuid.UUID) error {
	model, err := models.FindTemplatePart(ctx, r.db, partID.String())
	if err == sql.ErrNoRows {
		return template_part.ErrTemplatePartNotFound
	}
	if err != nil {
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}

func (r *TemplatePartRepository) DeleteByTemplateID(ctx context.Context, templateID uuid.UUID) error {
	_, err := models.TemplateParts(
		qm.Where("template_id = ?", templateID.String()),
	).DeleteAll(ctx, r.db)
	return err
}

func (r *TemplatePartRepository) List(ctx context.Context, limit, offset int) ([]*template_part.TemplatePart, error) {
	modelSlice, err := models.TemplateParts(
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	parts := make([]*template_part.TemplatePart, len(modelSlice))
	for i, model := range modelSlice {
		tp, err := toTemplatePartEntity(model)
		if err != nil {
			return nil, err
		}
		parts[i] = tp
	}
	return parts, nil
}
