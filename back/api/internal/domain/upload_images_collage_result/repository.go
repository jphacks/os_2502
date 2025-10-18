package upload_images_collage_result

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, relation *UploadImagesCollageResult) error
	FindByImageID(ctx context.Context, imageID uuid.UUID) ([]*UploadImagesCollageResult, error)
	FindByResultID(ctx context.Context, resultID uuid.UUID) ([]*UploadImagesCollageResult, error)
	Delete(ctx context.Context, imageID, resultID uuid.UUID) error
}
