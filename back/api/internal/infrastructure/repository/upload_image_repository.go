package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/upload_image"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/models"
)

type UploadImageRepositorySQLBoiler struct {
	db *sql.DB
}

func NewUploadImageRepositorySQLBoiler(db *sql.DB) upload_image.Repository {
	return &UploadImageRepositorySQLBoiler{db: db}
}

// Model to Entity conversion
func toUploadImageEntity(m *models.UploadImage) (*upload_image.UploadImage, error) {
	imageID, err := uuid.Parse(m.ImageID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(m.UserID)
	if err != nil {
		return nil, err
	}

	return upload_image.Reconstruct(
		imageID,
		m.FileURL,
		m.GroupID,
		userID,
		m.CollageDay,
		m.CreatedAt,
	)
}

// Entity to Model conversion
func toUploadImageModel(ui *upload_image.UploadImage) *models.UploadImage {
	return &models.UploadImage{
		ImageID:    ui.ImageID().String(),
		FileURL:    ui.FileURL(),
		GroupID:    ui.GroupID(),
		UserID:     ui.UserID().String(),
		CollageDay: ui.CollageDay(),
		CreatedAt:  ui.CreatedAt(),
	}
}

func (r *UploadImageRepositorySQLBoiler) Create(ctx context.Context, ui *upload_image.UploadImage) error {
	model := toUploadImageModel(ui)
	err := model.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		if db.IsDuplicateError(err) {
			return upload_image.ErrImageAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UploadImageRepositorySQLBoiler) FindByID(ctx context.Context, imageID uuid.UUID) (*upload_image.UploadImage, error) {
	model, err := models.FindUploadImage(ctx, r.db, imageID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, upload_image.ErrImageNotFound
		}
		return nil, err
	}
	return toUploadImageEntity(model)
}

func (r *UploadImageRepositorySQLBoiler) FindByGroupID(ctx context.Context, groupID string, limit, offset int) ([]*upload_image.UploadImage, error) {
	modelSlice, err := models.UploadImages(
		qm.Where("group_id = ?", groupID),
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	images := make([]*upload_image.UploadImage, len(modelSlice))
	for i, model := range modelSlice {
		img, err := toUploadImageEntity(model)
		if err != nil {
			return nil, err
		}
		images[i] = img
	}
	return images, nil
}

func (r *UploadImageRepositorySQLBoiler) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*upload_image.UploadImage, error) {
	modelSlice, err := models.UploadImages(
		qm.Where("user_id = ?", userID.String()),
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	images := make([]*upload_image.UploadImage, len(modelSlice))
	for i, model := range modelSlice {
		img, err := toUploadImageEntity(model)
		if err != nil {
			return nil, err
		}
		images[i] = img
	}
	return images, nil
}

func (r *UploadImageRepositorySQLBoiler) FindByGroupAndDate(ctx context.Context, groupID string, collageDay time.Time) ([]*upload_image.UploadImage, error) {
	modelSlice, err := models.UploadImages(
		qm.Where("group_id = ? AND collage_day = ?", groupID, collageDay.Format("2006-01-02")),
		qm.OrderBy("created_at ASC"),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	images := make([]*upload_image.UploadImage, len(modelSlice))
	for i, model := range modelSlice {
		img, err := toUploadImageEntity(model)
		if err != nil {
			return nil, err
		}
		images[i] = img
	}
	return images, nil
}

func (r *UploadImageRepositorySQLBoiler) FindByGroupUserAndDate(ctx context.Context, groupID string, userID uuid.UUID, collageDay time.Time) (*upload_image.UploadImage, error) {
	model, err := models.UploadImages(
		qm.Where("group_id = ? AND user_id = ? AND collage_day = ?", groupID, userID.String(), collageDay.Format("2006-01-02")),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, upload_image.ErrImageNotFound
		}
		return nil, err
	}
	return toUploadImageEntity(model)
}

func (r *UploadImageRepositorySQLBoiler) Delete(ctx context.Context, imageID uuid.UUID) error {
	model, err := models.FindUploadImage(ctx, r.db, imageID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return upload_image.ErrImageNotFound
		}
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}

func (r *UploadImageRepositorySQLBoiler) DeleteByGroupAndDate(ctx context.Context, groupID string, collageDay time.Time) (int, error) {
	count, err := models.UploadImages(
		qm.Where("group_id = ? AND collage_day = ?", groupID, collageDay.Format("2006-01-02")),
	).DeleteAll(ctx, r.db)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
