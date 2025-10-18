package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/upload_image"
)

type UploadImageUseCase struct {
	repo upload_image.Repository
}

func NewUploadImageUseCase(repo upload_image.Repository) *UploadImageUseCase {
	return &UploadImageUseCase{repo: repo}
}

// UploadImage uploads a new image
func (uc *UploadImageUseCase) UploadImage(ctx context.Context, fileURL, groupID string, userID uuid.UUID, collageDay time.Time) (*upload_image.UploadImage, error) {
	// 同じグループ、ユーザー、日付の画像が既に存在するかチェック
	existing, err := uc.repo.FindByGroupUserAndDate(ctx, groupID, userID, collageDay)
	if err == nil && existing != nil {
		// 既存の画像を削除して新しい画像に置き換える
		if err := uc.repo.Delete(ctx, existing.ImageID()); err != nil {
			return nil, err
		}
	}

	// 新規作成
	image, err := upload_image.NewUploadImage(fileURL, groupID, userID, collageDay)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, image); err != nil {
		return nil, err
	}

	return image, nil
}

// GetImage retrieves an image by ID
func (uc *UploadImageUseCase) GetImage(ctx context.Context, imageID uuid.UUID) (*upload_image.UploadImage, error) {
	return uc.repo.FindByID(ctx, imageID)
}

// GetImagesByGroup retrieves all images by group ID
func (uc *UploadImageUseCase) GetImagesByGroup(ctx context.Context, groupID string, limit, offset int) ([]*upload_image.UploadImage, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindByGroupID(ctx, groupID, limit, offset)
}

// DeleteImage deletes an image
func (uc *UploadImageUseCase) DeleteImage(ctx context.Context, imageID, userID uuid.UUID) error {
	// 画像を取得
	image, err := uc.repo.FindByID(ctx, imageID)
	if err != nil {
		return err
	}
	if image == nil {
		return upload_image.ErrImageNotFound
	}

	// ユーザー本人の画像かチェック
	if image.UserID() != userID {
		return upload_image.ErrNotAuthorized
	}

	return uc.repo.Delete(ctx, imageID)
}
