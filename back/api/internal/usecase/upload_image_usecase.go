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
	// 同じユーザーが同じ日に同じグループに画像をアップロード済みかチェック
	existing, err := uc.repo.FindByGroupUserAndDate(ctx, groupID, userID, collageDay)
	if err == nil && existing != nil {
		// 既存の画像を削除して新しい画像に置き換える
		if err := uc.repo.Delete(ctx, existing.ImageID()); err != nil {
			return nil, err
		}
	}

	image, err := upload_image.NewUploadImage(fileURL, groupID, userID, collageDay)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, image); err != nil {
		return nil, err
	}

	return image, nil
}

// GetImage gets an image by ID
func (uc *UploadImageUseCase) GetImage(ctx context.Context, imageID uuid.UUID) (*upload_image.UploadImage, error) {
	return uc.repo.FindByID(ctx, imageID)
}

// GetImagesByGroup gets all images by group ID
func (uc *UploadImageUseCase) GetImagesByGroup(ctx context.Context, groupID string, limit, offset int) ([]*upload_image.UploadImage, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindByGroupID(ctx, groupID, limit, offset)
}

// GetImagesByUser gets all images by user ID
func (uc *UploadImageUseCase) GetImagesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*upload_image.UploadImage, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindByUserID(ctx, userID, limit, offset)
}

// GetImagesByGroupAndDate gets all images by group ID and date
func (uc *UploadImageUseCase) GetImagesByGroupAndDate(ctx context.Context, groupID string, collageDay time.Time) ([]*upload_image.UploadImage, error) {
	return uc.repo.FindByGroupAndDate(ctx, groupID, collageDay)
}

// DeleteImage deletes an image
func (uc *UploadImageUseCase) DeleteImage(ctx context.Context, imageID uuid.UUID, userID uuid.UUID) error {
	image, err := uc.repo.FindByID(ctx, imageID)
	if err != nil {
		return err
	}

	// 画像の所有者であることを確認
	if image.UserID() != userID {
		return upload_image.ErrNotAuthorized
	}

	return uc.repo.Delete(ctx, imageID)
}

// DeleteImagesByGroupAndDate deletes all images by group ID and date
func (uc *UploadImageUseCase) DeleteImagesByGroupAndDate(ctx context.Context, groupID string, collageDay time.Time) (int, error) {
	return uc.repo.DeleteByGroupAndDate(ctx, groupID, collageDay)
}
