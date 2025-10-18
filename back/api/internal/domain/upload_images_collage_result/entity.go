package upload_images_collage_result

import (
	"time"

	"github.com/google/uuid"
)

type UploadImagesCollageResult struct {
	imageID    uuid.UUID
	resultID   uuid.UUID
	positionX  int
	positionY  int
	width      int
	height     int
	sortOrder  int
	createdAt  time.Time
}

func NewUploadImagesCollageResult(
	imageID, resultID uuid.UUID,
	positionX, positionY, width, height, sortOrder int,
) (*UploadImagesCollageResult, error) {
	if imageID == uuid.Nil {
		return nil, ErrInvalidImageID
	}
	if resultID == uuid.Nil {
		return nil, ErrInvalidResultID
	}
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidDimensions
	}

	return &UploadImagesCollageResult{
		imageID:   imageID,
		resultID:  resultID,
		positionX: positionX,
		positionY: positionY,
		width:     width,
		height:    height,
		sortOrder: sortOrder,
		createdAt: time.Now(),
	}, nil
}

func Reconstruct(
	imageID, resultID uuid.UUID,
	positionX, positionY, width, height, sortOrder int,
	createdAt time.Time,
) (*UploadImagesCollageResult, error) {
	if imageID == uuid.Nil {
		return nil, ErrInvalidImageID
	}
	if resultID == uuid.Nil {
		return nil, ErrInvalidResultID
	}
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidDimensions
	}

	return &UploadImagesCollageResult{
		imageID:   imageID,
		resultID:  resultID,
		positionX: positionX,
		positionY: positionY,
		width:     width,
		height:    height,
		sortOrder: sortOrder,
		createdAt: createdAt,
	}, nil
}

// Getters
func (uicr *UploadImagesCollageResult) ImageID() uuid.UUID {
	return uicr.imageID
}

func (uicr *UploadImagesCollageResult) ResultID() uuid.UUID {
	return uicr.resultID
}

func (uicr *UploadImagesCollageResult) PositionX() int {
	return uicr.positionX
}

func (uicr *UploadImagesCollageResult) PositionY() int {
	return uicr.positionY
}

func (uicr *UploadImagesCollageResult) Width() int {
	return uicr.width
}

func (uicr *UploadImagesCollageResult) Height() int {
	return uicr.height
}

func (uicr *UploadImagesCollageResult) SortOrder() int {
	return uicr.sortOrder
}

func (uicr *UploadImagesCollageResult) CreatedAt() time.Time {
	return uicr.createdAt
}

// UpdatePosition は画像の位置とサイズを更新
func (uicr *UploadImagesCollageResult) UpdatePosition(positionX, positionY, width, height int) error {
	if width <= 0 || height <= 0 {
		return ErrInvalidDimensions
	}
	uicr.positionX = positionX
	uicr.positionY = positionY
	uicr.width = width
	uicr.height = height
	return nil
}

// UpdateSortOrder は表示順序を更新
func (uicr *UploadImagesCollageResult) UpdateSortOrder(sortOrder int) {
	uicr.sortOrder = sortOrder
}
