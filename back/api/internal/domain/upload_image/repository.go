package upload_image

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	// Create creates a new upload image
	Create(ctx context.Context, image *UploadImage) error

	// FindByID finds an upload image by ID
	FindByID(ctx context.Context, imageID uuid.UUID) (*UploadImage, error)

	// FindByGroupID finds all upload images by group ID
	FindByGroupID(ctx context.Context, groupID string, limit, offset int) ([]*UploadImage, error)

	// FindByUserID finds all upload images by user ID
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*UploadImage, error)

	// FindByGroupAndDate finds all upload images by group ID and collage date
	FindByGroupAndDate(ctx context.Context, groupID string, collageDay time.Time) ([]*UploadImage, error)

	// FindByGroupUserAndDate finds an upload image by group ID, user ID and collage date
	FindByGroupUserAndDate(ctx context.Context, groupID string, userID uuid.UUID, collageDay time.Time) (*UploadImage, error)

	// Delete deletes an upload image
	Delete(ctx context.Context, imageID uuid.UUID) error

	// DeleteByGroupAndDate deletes all upload images by group ID and collage date
	DeleteByGroupAndDate(ctx context.Context, groupID string, collageDay time.Time) (int, error)
}
