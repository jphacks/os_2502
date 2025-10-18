package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/result_download"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db/models"
)

type ResultDownloadRepositorySQLBoiler struct {
	db *sql.DB
}

func NewResultDownloadRepositorySQLBoiler(db *sql.DB) result_download.Repository {
	return &ResultDownloadRepositorySQLBoiler{db: db}
}

// Model to Entity conversion
func toResultDownloadEntity(m *models.ResultDownload) (*result_download.ResultDownload, error) {
	resultID, err := uuid.Parse(m.ResultID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(m.UserID)
	if err != nil {
		return nil, err
	}

	return result_download.Reconstruct(
		resultID,
		userID,
		m.DownloadedAt,
	)
}

// Entity to Model conversion
func toResultDownloadModel(rd *result_download.ResultDownload) *models.ResultDownload {
	return &models.ResultDownload{
		ResultID:     rd.ResultID().String(),
		UserID:       rd.UserID().String(),
		DownloadedAt: rd.DownloadedAt(),
	}
}

func (r *ResultDownloadRepositorySQLBoiler) Create(ctx context.Context, rd *result_download.ResultDownload) error {
	model := toResultDownloadModel(rd)
	err := model.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		if db.IsDuplicateError(err) {
			return result_download.ErrDownloadAlreadyExists
		}
		return err
	}
	return nil
}

func (r *ResultDownloadRepositorySQLBoiler) FindByResultAndUser(ctx context.Context, resultID, userID uuid.UUID) (*result_download.ResultDownload, error) {
	model, err := models.ResultDownloads(
		qm.Where("result_id = ? AND user_id = ?", resultID.String(), userID.String()),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, result_download.ErrDownloadNotFound
		}
		return nil, err
	}
	return toResultDownloadEntity(model)
}

func (r *ResultDownloadRepositorySQLBoiler) FindByResultID(ctx context.Context, resultID uuid.UUID, limit, offset int) ([]*result_download.ResultDownload, error) {
	modelSlice, err := models.ResultDownloads(
		qm.Where("result_id = ?", resultID.String()),
		qm.OrderBy("downloaded_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	downloads := make([]*result_download.ResultDownload, len(modelSlice))
	for i, model := range modelSlice {
		rd, err := toResultDownloadEntity(model)
		if err != nil {
			return nil, err
		}
		downloads[i] = rd
	}
	return downloads, nil
}

func (r *ResultDownloadRepositorySQLBoiler) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*result_download.ResultDownload, error) {
	modelSlice, err := models.ResultDownloads(
		qm.Where("user_id = ?", userID.String()),
		qm.OrderBy("downloaded_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	downloads := make([]*result_download.ResultDownload, len(modelSlice))
	for i, model := range modelSlice {
		rd, err := toResultDownloadEntity(model)
		if err != nil {
			return nil, err
		}
		downloads[i] = rd
	}
	return downloads, nil
}

func (r *ResultDownloadRepositorySQLBoiler) CountByResultID(ctx context.Context, resultID uuid.UUID) (int, error) {
	count, err := models.ResultDownloads(
		qm.Where("result_id = ?", resultID.String()),
	).Count(ctx, r.db)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *ResultDownloadRepositorySQLBoiler) Delete(ctx context.Context, resultID, userID uuid.UUID) error {
	model, err := models.ResultDownloads(
		qm.Where("result_id = ? AND user_id = ?", resultID.String(), userID.String()),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return result_download.ErrDownloadNotFound
		}
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}
