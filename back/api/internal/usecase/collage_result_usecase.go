package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/collage_result"
)

type CollageResultUseCase struct {
	repo collage_result.Repository
}

func NewCollageResultUseCase(repo collage_result.Repository) *CollageResultUseCase {
	return &CollageResultUseCase{repo: repo}
}

// CreateResult creates a new collage result
func (uc *CollageResultUseCase) CreateResult(ctx context.Context, templateID uuid.UUID, groupID, fileURL string, targetUserNumber int) (*collage_result.CollageResult, error) {
	result, err := collage_result.NewCollageResult(templateID, groupID, fileURL, targetUserNumber)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetResult gets a result by ID
func (uc *CollageResultUseCase) GetResult(ctx context.Context, resultID uuid.UUID) (*collage_result.CollageResult, error) {
	return uc.repo.FindByID(ctx, resultID)
}

// GetResultsByGroup gets all results by group ID
func (uc *CollageResultUseCase) GetResultsByGroup(ctx context.Context, groupID string, limit, offset int) ([]*collage_result.CollageResult, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindByGroupID(ctx, groupID, limit, offset)
}

// GetUnnotifiedResults gets all unnotified results
func (uc *CollageResultUseCase) GetUnnotifiedResults(ctx context.Context, limit int) ([]*collage_result.CollageResult, error) {
	if limit <= 0 {
		limit = 100
	}
	return uc.repo.FindUnnotified(ctx, limit)
}

// MarkAsNotified marks a result as notified
func (uc *CollageResultUseCase) MarkAsNotified(ctx context.Context, resultID uuid.UUID) error {
	result, err := uc.repo.FindByID(ctx, resultID)
	if err != nil {
		return err
	}

	result.MarkAsNotified()
	return uc.repo.Update(ctx, result)
}

// DeleteResult deletes a result
func (uc *CollageResultUseCase) DeleteResult(ctx context.Context, resultID uuid.UUID) error {
	return uc.repo.Delete(ctx, resultID)
}
