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
	// 既に記録されているかチェック
	existing, err := uc.repo.FindByResultAndUser(ctx, resultID, userID)
	if err == nil && existing != nil {
		// 既に記録されている場合はそれを返す
		return existing, nil
	}

	// 新規作成
	download, err := result_download.NewResultDownload(resultID, userID)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, download); err != nil {
		return nil, err
	}

	return download, nil
}

// GetDownloadsByResult retrieves all downloads by result ID
func (uc *ResultDownloadUseCase) GetDownloadsByResult(ctx context.Context, resultID uuid.UUID, limit, offset int) ([]*result_download.ResultDownload, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindByResultID(ctx, resultID, limit, offset)
}

// GetDownloadCount retrieves the download count for a result
func (uc *ResultDownloadUseCase) GetDownloadCount(ctx context.Context, resultID uuid.UUID) (int, error) {
	return uc.repo.CountByResultID(ctx, resultID)
}
