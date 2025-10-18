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

func NewTemplatePart(templateID uuid.UUID, partNumber, positionX, positionY, width, height int) (*TemplatePart, error) {
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
		partID:     uuid.New(),
		templateID: templateID,
		partNumber: partNumber,
		positionX:  positionX,
		positionY:  positionY,
		width:      width,
		height:     height,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func Reconstruct(partID, templateID uuid.UUID, partNumber, positionX, positionY, width, height int,
	partName, description *string, createdAt, updatedAt time.Time) (*TemplatePart, error) {
	if partID == uuid.Nil {
		return nil, ErrInvalidPartID
	}
	if templateID == uuid.Nil {
		return nil, ErrInvalidTemplateID
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

// Setters
func (tp *TemplatePart) SetPartName(name string) {
	tp.partName = &name
	tp.updatedAt = time.Now()
}

func (tp *TemplatePart) SetDescription(desc string) {
	tp.description = &desc
	tp.updatedAt = time.Now()
}

func (tp *TemplatePart) UpdatePosition(x, y int) error {
	tp.positionX = x
	tp.positionY = y
	tp.updatedAt = time.Now()
	return nil
}

func (tp *TemplatePart) UpdateDimensions(width, height int) error {
	if width <= 0 || height <= 0 {
		return ErrInvalidDimensions
	}
	tp.width = width
	tp.height = height
	tp.updatedAt = time.Now()
	return nil
}
