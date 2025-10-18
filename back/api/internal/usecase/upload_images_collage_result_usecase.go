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

func (uc *UploadImagesCollageResultUseCase) LinkImageToResult(ctx context.Context, imageID, resultID uuid.UUID, positionX, positionY, width, height, sortOrder int) (*upload_images_collage_result.UploadImagesCollageResult, error) {
	// 画像とコラージュ結果を関連付け
	relation, err := upload_images_collage_result.NewUploadImagesCollageResult(imageID, resultID, positionX, positionY, width, height, sortOrder)
	if err != nil {
		return nil, err
	}

	// リポジトリに保存
	if err := uc.repo.Save(ctx, relation); err != nil {
		return nil, err
	}

	return relation, nil
}

func (uc *UploadImagesCollageResultUseCase) GetRelationsByImageID(ctx context.Context, imageID uuid.UUID) ([]*upload_images_collage_result.UploadImagesCollageResult, error) {
	return uc.repo.FindByImageID(ctx, imageID)
}

func (uc *UploadImagesCollageResultUseCase) GetRelationsByResultID(ctx context.Context, resultID uuid.UUID) ([]*upload_images_collage_result.UploadImagesCollageResult, error) {
	return uc.repo.FindByResultID(ctx, resultID)
}

func (uc *UploadImagesCollageResultUseCase) UnlinkImageFromResult(ctx context.Context, imageID, resultID uuid.UUID) error {
	return uc.repo.Delete(ctx, imageID, resultID)
}
