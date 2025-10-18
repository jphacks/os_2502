package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/device_token"
)

type DeviceTokenUseCase struct {
	repo device_token.Repository
}

func NewDeviceTokenUseCase(repo device_token.Repository) *DeviceTokenUseCase {
	return &DeviceTokenUseCase{repo: repo}
}

// RegisterDeviceToken registers a new device token or updates if it already exists
func (uc *DeviceTokenUseCase) RegisterDeviceToken(
	ctx context.Context,
	userID uuid.UUID,
	tokenString string,
	deviceType device_token.DeviceType,
	deviceName *string,
) (*device_token.DeviceToken, error) {
	// 既存のトークンを検索
	existingToken, err := uc.repo.FindByToken(ctx, tokenString)
	if err == nil && existingToken != nil {
		// トークンが既に存在する場合、アクティブ化して最終使用日時を更新
		existingToken.Activate()
		existingToken.UpdateLastUsedAt()
		if deviceName != nil {
			existingToken.UpdateDeviceName(deviceName)
		}
		if err := uc.repo.Update(ctx, existingToken); err != nil {
			return nil, err
		}
		return existingToken, nil
	}

	// 新規トークンを作成
	newToken, err := device_token.NewDeviceToken(userID, tokenString, deviceType, deviceName)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, newToken); err != nil {
		return nil, err
	}

	return newToken, nil
}

// GetDeviceToken gets a device token by ID
func (uc *DeviceTokenUseCase) GetDeviceToken(ctx context.Context, id uuid.UUID) (*device_token.DeviceToken, error) {
	return uc.repo.FindByID(ctx, id)
}

// GetUserDeviceTokens gets all device tokens for a user
func (uc *DeviceTokenUseCase) GetUserDeviceTokens(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*device_token.DeviceToken, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindByUserID(ctx, userID, limit, offset)
}

// GetActiveDeviceTokens gets all active device tokens for a user
func (uc *DeviceTokenUseCase) GetActiveDeviceTokens(ctx context.Context, userID uuid.UUID) ([]*device_token.DeviceToken, error) {
	return uc.repo.FindActiveByUserID(ctx, userID)
}

// DeactivateDeviceToken deactivates a device token
func (uc *DeviceTokenUseCase) DeactivateDeviceToken(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	token, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// トークンの所有者であることを確認
	if token.UserID() != userID {
		return device_token.ErrDeviceTokenNotFound
	}

	token.Deactivate()
	return uc.repo.Update(ctx, token)
}

// DeleteDeviceToken deletes a device token
func (uc *DeviceTokenUseCase) DeleteDeviceToken(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	token, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// トークンの所有者であることを確認
	if token.UserID() != userID {
		return device_token.ErrDeviceTokenNotFound
	}

	return uc.repo.Delete(ctx, id)
}

// UpdateLastUsed updates the last used timestamp of a device token
func (uc *DeviceTokenUseCase) UpdateLastUsed(ctx context.Context, tokenString string) error {
	token, err := uc.repo.FindByToken(ctx, tokenString)
	if err != nil {
		return err
	}

	token.UpdateLastUsedAt()
	return uc.repo.Update(ctx, token)
}

// CleanupOldTokens deactivates tokens that haven't been used for a specified number of days
func (uc *DeviceTokenUseCase) CleanupOldTokens(ctx context.Context, days int) (int, error) {
	if days <= 0 {
		days = 90 // デフォルト90日
	}
	return uc.repo.DeactivateOldTokens(ctx, days)
}
