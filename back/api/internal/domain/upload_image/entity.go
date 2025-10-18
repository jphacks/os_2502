package upload_image

import (
	"time"

	"github.com/google/uuid"
)

// UploadImage represents an uploaded image
type UploadImage struct {
	imageID    uuid.UUID
	fileURL    string
	groupID    string
	userID     uuid.UUID
	collageDay time.Time
	createdAt  time.Time
}

// NewUploadImage creates a new upload image
func NewUploadImage(fileURL, groupID string, userID uuid.UUID, collageDay time.Time) (*UploadImage, error) {
	if err := validateFileURL(fileURL); err != nil {
		return nil, err
	}

	if groupID == "" {
		return nil, ErrInvalidGroupID
	}

	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	return &UploadImage{
		imageID:    uuid.New(),
		fileURL:    fileURL,
		groupID:    groupID,
		userID:     userID,
		collageDay: collageDay,
		createdAt:  time.Now(),
	}, nil
}

// Reconstruct reconstructs an UploadImage from repository data
func Reconstruct(
	imageID uuid.UUID,
	fileURL string,
	groupID string,
	userID uuid.UUID,
	collageDay time.Time,
	createdAt time.Time,
) (*UploadImage, error) {
	return &UploadImage{
		imageID:    imageID,
		fileURL:    fileURL,
		groupID:    groupID,
		userID:     userID,
		collageDay: collageDay,
		createdAt:  createdAt,
	}, nil
}

// Getters
func (ui *UploadImage) ImageID() uuid.UUID {
	return ui.imageID
}

func (ui *UploadImage) FileURL() string {
	return ui.fileURL
}

func (ui *UploadImage) GroupID() string {
	return ui.groupID
}

func (ui *UploadImage) UserID() uuid.UUID {
	return ui.userID
}

func (ui *UploadImage) CollageDay() time.Time {
	return ui.collageDay
}

func (ui *UploadImage) CreatedAt() time.Time {
	return ui.createdAt
}

// Validation functions
func validateFileURL(fileURL string) error {
	if fileURL == "" {
		return ErrInvalidFileURL
	}
	if len(fileURL) > 500 {
		return ErrInvalidFileURL
	}
	return nil
}
