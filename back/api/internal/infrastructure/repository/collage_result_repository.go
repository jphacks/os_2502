package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/noonyuu/collage/api/internal/domain/collage_result"
	"github.com/noonyuu/collage/api/internal/infrastructure/db"
	"github.com/noonyuu/collage/api/internal/infrastructure/models"
)

type CollageResultRepositorySQLBoiler struct {
	db *sql.DB
}

func NewCollageResultRepositorySQLBoiler(db *sql.DB) collage_result.Repository {
	return &CollageResultRepositorySQLBoiler{db: db}
}

// Model to Entity conversion
func toCollageResultEntity(m *models.CollageResult) (*collage_result.CollageResult, error) {
	resultID, err := uuid.Parse(m.ResultID)
	if err != nil {
		return nil, err
	}

	templateID, err := uuid.Parse(m.TemplateID)
	if err != nil {
		return nil, err
	}

	return collage_result.Reconstruct(
		resultID,
		templateID,
		m.GroupID,
		m.FileURL,
		m.TargetUserNumber,
		m.IsNotification,
		m.CreatedAt,
	)
}

// Entity to Model conversion
func toCollageResultModel(cr *collage_result.CollageResult) *models.CollageResult {
	return &models.CollageResult{
		ResultID:         cr.ResultID().String(),
		TemplateID:       cr.TemplateID().String(),
		GroupID:          cr.GroupID(),
		FileURL:          cr.FileURL(),
		TargetUserNumber: cr.TargetUserNumber(),
		IsNotification:   cr.IsNotification(),
		CreatedAt:        cr.CreatedAt(),
	}
}

func (r *CollageResultRepositorySQLBoiler) Create(ctx context.Context, cr *collage_result.CollageResult) error {
	model := toCollageResultModel(cr)
	err := model.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		if db.IsDuplicateError(err) {
			return collage_result.ErrResultAlreadyExists
		}
		return err
	}
	return nil
}

func (r *CollageResultRepositorySQLBoiler) FindByID(ctx context.Context, resultID uuid.UUID) (*collage_result.CollageResult, error) {
	model, err := models.FindCollageResult(ctx, r.db, resultID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, collage_result.ErrResultNotFound
		}
		return nil, err
	}
	return toCollageResultEntity(model)
}

func (r *CollageResultRepositorySQLBoiler) FindByGroupID(ctx context.Context, groupID string, limit, offset int) ([]*collage_result.CollageResult, error) {
	modelSlice, err := models.CollageResults(
		qm.Where("group_id = ?", groupID),
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	results := make([]*collage_result.CollageResult, len(modelSlice))
	for i, model := range modelSlice {
		cr, err := toCollageResultEntity(model)
		if err != nil {
			return nil, err
		}
		results[i] = cr
	}
	return results, nil
}

func (r *CollageResultRepositorySQLBoiler) FindUnnotified(ctx context.Context, limit int) ([]*collage_result.CollageResult, error) {
	modelSlice, err := models.CollageResults(
		qm.Where("is_notification = ?", false),
		qm.OrderBy("created_at ASC"),
		qm.Limit(limit),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	results := make([]*collage_result.CollageResult, len(modelSlice))
	for i, model := range modelSlice {
		cr, err := toCollageResultEntity(model)
		if err != nil {
			return nil, err
		}
		results[i] = cr
	}
	return results, nil
}

func (r *CollageResultRepositorySQLBoiler) Update(ctx context.Context, cr *collage_result.CollageResult) error {
	model, err := models.FindCollageResult(ctx, r.db, cr.ResultID().String())
	if err != nil {
		if err == sql.ErrNoRows {
			return collage_result.ErrResultNotFound
		}
		return err
	}

	model.IsNotification = cr.IsNotification()

	_, err = model.Update(ctx, r.db, boil.Whitelist(
		models.CollageResultColumns.IsNotification,
	))
	return err
}

func (r *CollageResultRepositorySQLBoiler) Delete(ctx context.Context, resultID uuid.UUID) error {
	model, err := models.FindCollageResult(ctx, r.db, resultID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return collage_result.ErrResultNotFound
		}
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}
