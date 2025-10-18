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

// RegisterDeviceToken registers or updates a device token
func (uc *DeviceTokenUseCase) RegisterDeviceToken(ctx context.Context, userID uuid.UUID, deviceTokenStr string, deviceType device_token.DeviceType, deviceName *string) (*device_token.DeviceToken, error) {
	// 既存のトークンをチェック
	existing, err := uc.repo.FindByToken(ctx, deviceTokenStr)
	if err == nil && existing != nil {
		// 既に存在する場合は、最終使用時刻を更新してアクティブにする
		existing.UpdateLastUsedAt()
		if !existing.IsActive() {
			existing.Activate()
		}
		if err := uc.repo.Update(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// 新規作成
	dt, err := device_token.NewDeviceToken(userID, deviceTokenStr, deviceType, deviceName)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, dt); err != nil {
		return nil, err
	}

	return dt, nil
}

// GetDeviceToken retrieves a device token by ID
func (uc *DeviceTokenUseCase) GetDeviceToken(ctx context.Context, id uuid.UUID) (*device_token.DeviceToken, error) {
	return uc.repo.FindByID(ctx, id)
}

// GetUserDeviceTokens retrieves all device tokens for a user
func (uc *DeviceTokenUseCase) GetUserDeviceTokens(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*device_token.DeviceToken, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindByUserID(ctx, userID, limit, offset)
}

// DeactivateDeviceToken deactivates a device token
func (uc *DeviceTokenUseCase) DeactivateDeviceToken(ctx context.Context, id, userID uuid.UUID) error {
	// デバイストークンを取得
	dt, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if dt == nil {
		return device_token.ErrDeviceTokenNotFound
	}

	// ユーザー本人のトークンかチェック
	if dt.UserID() != userID {
		return device_token.ErrDeviceTokenNotFound
	}

	// 無効化
	dt.Deactivate()

	return uc.repo.Update(ctx, dt)
}

// DeleteDeviceToken deletes a device token
func (uc *DeviceTokenUseCase) DeleteDeviceToken(ctx context.Context, id, userID uuid.UUID) error {
	// デバイストークンを取得
	dt, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if dt == nil {
		return device_token.ErrDeviceTokenNotFound
	}

	// ユーザー本人のトークンかチェック
	if dt.UserID() != userID {
		return device_token.ErrDeviceTokenNotFound
	}

	// 削除
	return uc.repo.Delete(ctx, id)
}
