package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/upload_images_collage_result"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db/models"
)

type UploadImagesCollageResultRepositorySQLBoiler struct {
	db *sql.DB
}

func NewUploadImagesCollageResultRepositorySQLBoiler(db *sql.DB) upload_images_collage_result.Repository {
	return &UploadImagesCollageResultRepositorySQLBoiler{db: db}
}

func toUploadImagesCollageResultModel(uicr *upload_images_collage_result.UploadImagesCollageResult) *models.UploadImagesCollageResult {
	return &models.UploadImagesCollageResult{
		ImageID:   uicr.ImageID().String(),
		ResultID:  uicr.ResultID().String(),
		PositionX: uicr.PositionX(),
		PositionY: uicr.PositionY(),
		Width:     uicr.Width(),
		Height:    uicr.Height(),
		SortOrder: uicr.SortOrder(),
		CreatedAt: uicr.CreatedAt(),
	}
}

func toUploadImagesCollageResultEntity(m *models.UploadImagesCollageResult) (*upload_images_collage_result.UploadImagesCollageResult, error) {
	imageID, err := uuid.Parse(m.ImageID)
	if err != nil {
		return nil, err
	}
	resultID, err := uuid.Parse(m.ResultID)
	if err != nil {
		return nil, err
	}

	return upload_images_collage_result.Reconstruct(
		imageID,
		resultID,
		m.PositionX,
		m.PositionY,
		m.Width,
		m.Height,
		m.SortOrder,
		m.CreatedAt,
	)
}

func (r *UploadImagesCollageResultRepositorySQLBoiler) Save(ctx context.Context, relation *upload_images_collage_result.UploadImagesCollageResult) error {
	model := toUploadImagesCollageResultModel(relation)
	return model.Upsert(ctx, r.db, true, []string{"image_id", "result_id"}, boil.Infer(), boil.Infer())
}

func (r *UploadImagesCollageResultRepositorySQLBoiler) FindByImageID(ctx context.Context, imageID uuid.UUID) ([]*upload_images_collage_result.UploadImagesCollageResult, error) {
	modelSlice, err := models.UploadImagesCollageResults(
		qm.Where("image_id = ?", imageID.String()),
		qm.OrderBy("sort_order ASC"),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	entities := make([]*upload_images_collage_result.UploadImagesCollageResult, len(modelSlice))
	for i, model := range modelSlice {
		entity, err := toUploadImagesCollageResultEntity(model)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}

func (r *UploadImagesCollageResultRepositorySQLBoiler) FindByResultID(ctx context.Context, resultID uuid.UUID) ([]*upload_images_collage_result.UploadImagesCollageResult, error) {
	modelSlice, err := models.UploadImagesCollageResults(
		qm.Where("result_id = ?", resultID.String()),
		qm.OrderBy("sort_order ASC"),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	entities := make([]*upload_images_collage_result.UploadImagesCollageResult, len(modelSlice))
	for i, model := range modelSlice {
		entity, err := toUploadImagesCollageResultEntity(model)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}

func (r *UploadImagesCollageResultRepositorySQLBoiler) Delete(ctx context.Context, imageID, resultID uuid.UUID) error {
	model, err := models.FindUploadImagesCollageResult(ctx, r.db, imageID.String(), resultID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return upload_images_collage_result.ErrNotFound
		}
		return err
	}
	_, err = model.Delete(ctx, r.db)
	return err
}
