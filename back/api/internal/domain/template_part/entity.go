package template_part

import (
	"time"

	"github.com/google/uuid"
)

type TemplatePart struct {
	partID      uuid.UUID
	templateID  uuid.UUID
	partNumber  int
	partName    *string
	positionX   int
	positionY   int
	width       int
	height      int
	description *string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewTemplatePart(
	templateID uuid.UUID,
	partNumber, positionX, positionY, width, height int,
	partName, description *string,
) (*TemplatePart, error) {
	if templateID == uuid.Nil {
		return nil, ErrInvalidTemplateID
	}
	if partNumber < 1 {
		return nil, ErrInvalidPartNumber
	}
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidDimensions
	}

	now := time.Now()
	return &TemplatePart{
		partID:      uuid.New(),
		templateID:  templateID,
		partNumber:  partNumber,
		partName:    partName,
		positionX:   positionX,
		positionY:   positionY,
		width:       width,
		height:      height,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func Reconstruct(
	partID, templateID uuid.UUID,
	partNumber, positionX, positionY, width, height int,
	partName, description *string,
	createdAt, updatedAt time.Time,
) (*TemplatePart, error) {
	if partID == uuid.Nil {
		return nil, ErrInvalidPartID
	}
	if templateID == uuid.Nil {
		return nil, ErrInvalidTemplateID
	}
	if partNumber < 1 {
		return nil, ErrInvalidPartNumber
	}
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidDimensions
	}

	return &TemplatePart{
		partID:      partID,
		templateID:  templateID,
		partNumber:  partNumber,
		partName:    partName,
		positionX:   positionX,
		positionY:   positionY,
		width:       width,
		height:      height,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

// Getters
func (tp *TemplatePart) PartID() uuid.UUID {
	return tp.partID
}

func (tp *TemplatePart) TemplateID() uuid.UUID {
	return tp.templateID
}

func (tp *TemplatePart) PartNumber() int {
	return tp.partNumber
}

func (tp *TemplatePart) PartName() *string {
	return tp.partName
}

func (tp *TemplatePart) PositionX() int {
	return tp.positionX
}

func (tp *TemplatePart) PositionY() int {
	return tp.positionY
}

func (tp *TemplatePart) Width() int {
	return tp.width
}

func (tp *TemplatePart) Height() int {
	return tp.height
}

func (tp *TemplatePart) Description() *string {
	return tp.description
}

func (tp *TemplatePart) CreatedAt() time.Time {
	return tp.createdAt
}

func (tp *TemplatePart) UpdatedAt() time.Time {
	return tp.updatedAt
}

// UpdatePosition はパーツの位置とサイズを更新
func (tp *TemplatePart) UpdatePosition(positionX, positionY, width, height int) error {
	if width <= 0 || height <= 0 {
		return ErrInvalidDimensions
	}
	tp.positionX = positionX
	tp.positionY = positionY
	tp.width = width
	tp.height = height
	tp.updatedAt = time.Now()
	return nil
}

// UpdatePartName はパーツ名を更新
func (tp *TemplatePart) UpdatePartName(partName *string) {
	tp.partName = partName
	tp.updatedAt = time.Now()
}

// UpdateDescription は説明を更新
func (tp *TemplatePart) UpdateDescription(description *string) {
	tp.description = description
	tp.updatedAt = time.Now()
}
