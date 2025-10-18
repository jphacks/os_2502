package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/result_download"
)

type ResultDownloadUseCase struct {
	repo result_download.Repository
}

func NewResultDownloadUseCase(repo result_download.Repository) *ResultDownloadUseCase {
	return &ResultDownloadUseCase{repo: repo}
}

// RecordDownload records a download
func (uc *ResultDownloadUseCase) RecordDownload(ctx context.Context, resultID, userID uuid.UUID) (*result_download.ResultDownload, error) {
	// 既にダウンロード履歴があるかチェック
	existing, err := uc.repo.FindByResultAndUser(ctx, resultID, userID)
	if err == nil && existing != nil {
		// 既に存在する場合はそのまま返す
		return existing, nil
	}

	download, err := result_download.NewResultDownload(resultID, userID)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, download); err != nil {
		return nil, err
	}

	return download, nil
}

// GetDownloadsByResult gets all downloads by result ID
func (uc *ResultDownloadUseCase) GetDownloadsByResult(ctx context.Context, resultID uuid.UUID, limit, offset int) ([]*result_download.ResultDownload, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindByResultID(ctx, resultID, limit, offset)
}

// GetDownloadsByUser gets all downloads by user ID
func (uc *ResultDownloadUseCase) GetDownloadsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*result_download.ResultDownload, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindByUserID(ctx, userID, limit, offset)
}

// GetDownloadCount gets download count by result ID
func (uc *ResultDownloadUseCase) GetDownloadCount(ctx context.Context, resultID uuid.UUID) (int, error) {
	return uc.repo.CountByResultID(ctx, resultID)
}

// CheckDownloaded checks if a user has downloaded a result
func (uc *ResultDownloadUseCase) CheckDownloaded(ctx context.Context, resultID, userID uuid.UUID) (bool, error) {
	_, err := uc.repo.FindByResultAndUser(ctx, resultID, userID)
	if err == result_download.ErrDownloadNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
