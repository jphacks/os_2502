package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/upload_images_collage_result"
)

type UploadImagesCollageResultUseCase struct {
	repo upload_images_collage_result.Repository
}

func NewUploadImagesCollageResultUseCase(repo upload_images_collage_result.Repository) *UploadImagesCollageResultUseCase {
	return &UploadImagesCollageResultUseCase{repo: repo}
}

func (uc *UploadImagesCollageResultUseCase) CreateUploadImagesCollageResult(
	ctx context.Context,
	imageID, resultID uuid.UUID,
	positionX, positionY, width, height, sortOrder int,
) (*upload_images_collage_result.UploadImagesCollageResult, error) {
	// 同じ画像IDとコラージュ結果IDの組み合わせが既に存在するかチェック
	existing, err := uc.repo.FindByImageIDAndResultID(ctx, imageID, resultID)
	if err == nil && existing != nil {
		return nil, upload_images_collage_result.ErrUploadImagesCollageResultAlreadyExists
	}

	uicr, err := upload_images_collage_result.NewUploadImagesCollageResult(imageID, resultID, positionX, positionY, width, height, sortOrder)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, uicr); err != nil {
		return nil, err
	}

	return uicr, nil
}

func (uc *UploadImagesCollageResultUseCase) GetUploadImagesCollageResultByImageIDAndResultID(ctx context.Context, imageID, resultID uuid.UUID) (*upload_images_collage_result.UploadImagesCollageResult, error) {
	return uc.repo.FindByImageIDAndResultID(ctx, imageID, resultID)
}

func (uc *UploadImagesCollageResultUseCase) GetUploadImagesCollageResultsByImageID(ctx context.Context, imageID uuid.UUID) ([]*upload_images_collage_result.UploadImagesCollageResult, error) {
	return uc.repo.FindByImageID(ctx, imageID)
}

func (uc *UploadImagesCollageResultUseCase) GetUploadImagesCollageResultsByResultID(ctx context.Context, resultID uuid.UUID) ([]*upload_images_collage_result.UploadImagesCollageResult, error) {
	return uc.repo.FindByResultID(ctx, resultID)
}

func (uc *UploadImagesCollageResultUseCase) UpdateUploadImagesCollageResultPosition(
	ctx context.Context,
	imageID, resultID uuid.UUID,
	positionX, positionY, width, height int,
) (*upload_images_collage_result.UploadImagesCollageResult, error) {
	uicr, err := uc.repo.FindByImageIDAndResultID(ctx, imageID, resultID)
	if err != nil {
		return nil, err
	}

	if err := uicr.UpdatePosition(positionX, positionY, width, height); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, uicr); err != nil {
		return nil, err
	}

	return uicr, nil
}

func (uc *UploadImagesCollageResultUseCase) UpdateUploadImagesCollageResultSortOrder(
	ctx context.Context,
	imageID, resultID uuid.UUID,
	sortOrder int,
) (*upload_images_collage_result.UploadImagesCollageResult, error) {
	uicr, err := uc.repo.FindByImageIDAndResultID(ctx, imageID, resultID)
	if err != nil {
		return nil, err
	}

	uicr.UpdateSortOrder(sortOrder)

	if err := uc.repo.Update(ctx, uicr); err != nil {
		return nil, err
	}

	return uicr, nil
}

func (uc *UploadImagesCollageResultUseCase) DeleteUploadImagesCollageResult(ctx context.Context, imageID, resultID uuid.UUID) error {
	return uc.repo.Delete(ctx, imageID, resultID)
}

func (uc *UploadImagesCollageResultUseCase) DeleteUploadImagesCollageResultsByResultID(ctx context.Context, resultID uuid.UUID) error {
	return uc.repo.DeleteByResultID(ctx, resultID)
}

func (uc *UploadImagesCollageResultUseCase) ListUploadImagesCollageResults(ctx context.Context, limit, offset int) ([]*upload_images_collage_result.UploadImagesCollageResult, error) {
	return uc.repo.List(ctx, limit, offset)
}
