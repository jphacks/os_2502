package device_token

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Create creates a new device token
	Create(ctx context.Context, deviceToken *DeviceToken) error

	// FindByID finds a device token by ID
	FindByID(ctx context.Context, id uuid.UUID) (*DeviceToken, error)

	// FindByToken finds a device token by token string
	FindByToken(ctx context.Context, token string) (*DeviceToken, error)

	// FindByUserID finds all device tokens for a user
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*DeviceToken, error)

	// FindActiveByUserID finds all active device tokens for a user
	FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*DeviceToken, error)

	// Update updates a device token
	Update(ctx context.Context, deviceToken *DeviceToken) error

	// Delete deletes a device token
	Delete(ctx context.Context, id uuid.UUID) error

	// DeactivateOldTokens deactivates device tokens that haven't been used for a certain period
	DeactivateOldTokens(ctx context.Context, days int) (int, error)
}
