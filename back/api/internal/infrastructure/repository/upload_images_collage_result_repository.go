package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/upload_images_collage_result"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/models"
)

type UploadImagesCollageResultRepository struct {
	db *sql.DB
}

func NewUploadImagesCollageResultRepository(db *sql.DB) upload_images_collage_result.Repository {
	return &UploadImagesCollageResultRepository{db: db}
}

// Model to Entity conversion
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

// Entity to Model conversion
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

func (r *UploadImagesCollageResultRepository) Create(ctx context.Context, relation *upload_images_collage_result.UploadImagesCollageResult) error {
	model := toUploadImagesCollageResultModel(relation)
	return model.Insert(ctx, r.db, boil.Infer())
}

func (r *UploadImagesCollageResultRepository) FindByImageIDAndResultID(ctx context.Context, imageID, resultID uuid.UUID) (*upload_images_collage_result.UploadImagesCollageResult, error) {
	model, err := models.UploadImagesCollageResults(
		qm.Where("image_id = ? AND result_id = ?", imageID.String(), resultID.String()),
	).One(ctx, r.db)
	if err == sql.ErrNoRows {
		return nil, upload_images_collage_result.ErrUploadImagesCollageResultNotFound
	}
	if err != nil {
		return nil, err
	}
	return toUploadImagesCollageResultEntity(model)
}

func (r *UploadImagesCollageResultRepository) FindByImageID(ctx context.Context, imageID uuid.UUID) ([]*upload_images_collage_result.UploadImagesCollageResult, error) {
	modelSlice, err := models.UploadImagesCollageResults(
		qm.Where("image_id = ?", imageID.String()),
		qm.OrderBy("sort_order"),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	relations := make([]*upload_images_collage_result.UploadImagesCollageResult, len(modelSlice))
	for i, model := range modelSlice {
		uicr, err := toUploadImagesCollageResultEntity(model)
		if err != nil {
			return nil, err
		}
		relations[i] = uicr
	}
	return relations, nil
}

func (r *UploadImagesCollageResultRepository) FindByResultID(ctx context.Context, resultID uuid.UUID) ([]*upload_images_collage_result.UploadImagesCollageResult, error) {
	modelSlice, err := models.UploadImagesCollageResults(
		qm.Where("result_id = ?", resultID.String()),
		qm.OrderBy("sort_order"),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	relations := make([]*upload_images_collage_result.UploadImagesCollageResult, len(modelSlice))
	for i, model := range modelSlice {
		uicr, err := toUploadImagesCollageResultEntity(model)
		if err != nil {
			return nil, err
		}
		relations[i] = uicr
	}
	return relations, nil
}

func (r *UploadImagesCollageResultRepository) Update(ctx context.Context, relation *upload_images_collage_result.UploadImagesCollageResult) error {
	model, err := models.UploadImagesCollageResults(
		qm.Where("image_id = ? AND result_id = ?", relation.ImageID().String(), relation.ResultID().String()),
	).One(ctx, r.db)
	if err == sql.ErrNoRows {
		return upload_images_collage_result.ErrUploadImagesCollageResultNotFound
	}
	if err != nil {
		return err
	}

	model.PositionX = relation.PositionX()
	model.PositionY = relation.PositionY()
	model.Width = relation.Width()
	model.Height = relation.Height()
	model.SortOrder = relation.SortOrder()

	_, err = model.Update(ctx, r.db, boil.Whitelist(
		models.UploadImagesCollageResultColumns.PositionX,
		models.UploadImagesCollageResultColumns.PositionY,
		models.UploadImagesCollageResultColumns.Width,
		models.UploadImagesCollageResultColumns.Height,
		models.UploadImagesCollageResultColumns.SortOrder,
	))
	return err
}

func (r *UploadImagesCollageResultRepository) Delete(ctx context.Context, imageID, resultID uuid.UUID) error {
	count, err := models.UploadImagesCollageResults(
		qm.Where("image_id = ? AND result_id = ?", imageID.String(), resultID.String()),
	).DeleteAll(ctx, r.db)
	if err != nil {
		return err
	}
	if count == 0 {
		return upload_images_collage_result.ErrUploadImagesCollageResultNotFound
	}
	return nil
}

func (r *UploadImagesCollageResultRepository) DeleteByResultID(ctx context.Context, resultID uuid.UUID) error {
	_, err := models.UploadImagesCollageResults(
		qm.Where("result_id = ?", resultID.String()),
	).DeleteAll(ctx, r.db)
	return err
}

func (r *UploadImagesCollageResultRepository) List(ctx context.Context, limit, offset int) ([]*upload_images_collage_result.UploadImagesCollageResult, error) {
	modelSlice, err := models.UploadImagesCollageResults(
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	relations := make([]*upload_images_collage_result.UploadImagesCollageResult, len(modelSlice))
	for i, model := range modelSlice {
		uicr, err := toUploadImagesCollageResultEntity(model)
		if err != nil {
			return nil, err
		}
		relations[i] = uicr
	}
	return relations, nil
}
