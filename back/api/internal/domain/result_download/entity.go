package result_download

import (
	"time"

	"github.com/google/uuid"
)

// ResultDownload represents a result download record
type ResultDownload struct {
	resultID     uuid.UUID
	userID       uuid.UUID
	downloadedAt time.Time
}

// NewResultDownload creates a new result download record
func NewResultDownload(resultID, userID uuid.UUID) (*ResultDownload, error) {
	if resultID == uuid.Nil {
		return nil, ErrInvalidResultID
	}

	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	return &ResultDownload{
		resultID:     resultID,
		userID:       userID,
		downloadedAt: time.Now(),
	}, nil
}

// Reconstruct reconstructs a ResultDownload from repository data
func Reconstruct(
	resultID uuid.UUID,
	userID uuid.UUID,
	downloadedAt time.Time,
) (*ResultDownload, error) {
	return &ResultDownload{
		resultID:     resultID,
		userID:       userID,
		downloadedAt: downloadedAt,
	}, nil
}

// Getters
func (rd *ResultDownload) ResultID() uuid.UUID {
	return rd.resultID
}

func (rd *ResultDownload) UserID() uuid.UUID {
	return rd.userID
}

func (rd *ResultDownload) DownloadedAt() time.Time {
	return rd.downloadedAt
}
