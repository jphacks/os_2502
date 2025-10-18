package result_download

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Create creates a new result download record
	Create(ctx context.Context, download *ResultDownload) error

	// FindByResultAndUser finds a download record by result ID and user ID
	FindByResultAndUser(ctx context.Context, resultID, userID uuid.UUID) (*ResultDownload, error)

	// FindByResultID finds all download records by result ID
	FindByResultID(ctx context.Context, resultID uuid.UUID, limit, offset int) ([]*ResultDownload, error)

	// FindByUserID finds all download records by user ID
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ResultDownload, error)

	// CountByResultID counts download records by result ID
	CountByResultID(ctx context.Context, resultID uuid.UUID) (int, error)

	// Delete deletes a download record
	Delete(ctx context.Context, resultID, userID uuid.UUID) error
}
