package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/collage_template"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db/models"
)

type CollageTemplateRepositorySQLBoiler struct {
	db *sql.DB
}

func NewCollageTemplateRepositorySQLBoiler(db *sql.DB) collage_template.Repository {
	return &CollageTemplateRepositorySQLBoiler{db: db}
}

// Model to Entity conversion
func toCollageTemplateEntity(m *models.CollagesTemplate) (*collage_template.CollageTemplate, error) {
	templateID, err := uuid.Parse(m.TemplateID)
	if err != nil {
		return nil, err
	}

	return collage_template.Reconstruct(
		templateID,
		m.Name,
		m.FilePath,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// Entity to Model conversion
func toCollageTemplateModel(ct *collage_template.CollageTemplate) *models.CollagesTemplate {
	return &models.CollagesTemplate{
		TemplateID: ct.TemplateID().String(),
		Name:       ct.Name(),
		FilePath:   ct.FilePath(),
		CreatedAt:  ct.CreatedAt(),
		UpdatedAt:  ct.UpdatedAt(),
	}
}

func (r *CollageTemplateRepositorySQLBoiler) Create(ctx context.Context, ct *collage_template.CollageTemplate) error {
	model := toCollageTemplateModel(ct)
	err := model.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		if db.IsDuplicateError(err) {
			return collage_template.ErrTemplateAlreadyExists
		}
		return err
	}
	return nil
}

func (r *CollageTemplateRepositorySQLBoiler) FindByID(ctx context.Context, templateID uuid.UUID) (*collage_template.CollageTemplate, error) {
	model, err := models.FindCollagesTemplate(ctx, r.db, templateID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, collage_template.ErrTemplateNotFound
		}
		return nil, err
	}
	return toCollageTemplateEntity(model)
}

func (r *CollageTemplateRepositorySQLBoiler) FindByName(ctx context.Context, name string) (*collage_template.CollageTemplate, error) {
	model, err := models.CollagesTemplates(
		qm.Where("name = ?", name),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, collage_template.ErrTemplateNotFound
		}
		return nil, err
	}
	return toCollageTemplateEntity(model)
}

func (r *CollageTemplateRepositorySQLBoiler) List(ctx context.Context, limit, offset int) ([]*collage_template.CollageTemplate, error) {
	modelSlice, err := models.CollagesTemplates(
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	templates := make([]*collage_template.CollageTemplate, len(modelSlice))
	for i, model := range modelSlice {
		ct, err := toCollageTemplateEntity(model)
		if err != nil {
			return nil, err
		}
		templates[i] = ct
	}
	return templates, nil
}

func (r *CollageTemplateRepositorySQLBoiler) Update(ctx context.Context, ct *collage_template.CollageTemplate) error {
	model, err := models.FindCollagesTemplate(ctx, r.db, ct.TemplateID().String())
	if err != nil {
		if err == sql.ErrNoRows {
			return collage_template.ErrTemplateNotFound
		}
		return err
	}

	model.Name = ct.Name()
	model.FilePath = ct.FilePath()
	model.UpdatedAt = ct.UpdatedAt()

	_, err = model.Update(ctx, r.db, boil.Whitelist(
		models.CollagesTemplateColumns.Name,
		models.CollagesTemplateColumns.FilePath,
		models.CollagesTemplateColumns.UpdatedAt,
	))
	return err
}

func (r *CollageTemplateRepositorySQLBoiler) Delete(ctx context.Context, templateID uuid.UUID) error {
	model, err := models.FindCollagesTemplate(ctx, r.db, templateID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return collage_template.ErrTemplateNotFound
		}
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}
