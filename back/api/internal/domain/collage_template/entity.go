package collage_template

import (
	"time"

	"github.com/google/uuid"
)

// CollageTemplate represents a collage template
type CollageTemplate struct {
	templateID uuid.UUID
	name       string
	filePath   string
	createdAt  time.Time
	updatedAt  time.Time
}

// NewCollageTemplate creates a new collage template
func NewCollageTemplate(name, filePath string) (*CollageTemplate, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	if err := validateFilePath(filePath); err != nil {
		return nil, err
	}

	now := time.Now()
	return &CollageTemplate{
		templateID: uuid.New(),
		name:       name,
		filePath:   filePath,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

// Reconstruct reconstructs a CollageTemplate from repository data
func Reconstruct(
	templateID uuid.UUID,
	name string,
	filePath string,
	createdAt time.Time,
	updatedAt time.Time,
) (*CollageTemplate, error) {
	return &CollageTemplate{
		templateID: templateID,
		name:       name,
		filePath:   filePath,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}, nil
}

// Getters
func (ct *CollageTemplate) TemplateID() uuid.UUID {
	return ct.templateID
}

func (ct *CollageTemplate) Name() string {
	return ct.name
}

func (ct *CollageTemplate) FilePath() string {
	return ct.filePath
}

func (ct *CollageTemplate) CreatedAt() time.Time {
	return ct.createdAt
}

func (ct *CollageTemplate) UpdatedAt() time.Time {
	return ct.updatedAt
}

// UpdateName updates the template name
func (ct *CollageTemplate) UpdateName(name string) error {
	if err := validateName(name); err != nil {
		return err
	}
	ct.name = name
	ct.updatedAt = time.Now()
	return nil
}

// UpdateFilePath updates the template file path
func (ct *CollageTemplate) UpdateFilePath(filePath string) error {
	if err := validateFilePath(filePath); err != nil {
		return err
	}
	ct.filePath = filePath
	ct.updatedAt = time.Now()
	return nil
}

// Validation functions
func validateName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	if len(name) > 100 {
		return ErrInvalidName
	}
	return nil
}

func validateFilePath(filePath string) error {
	if filePath == "" {
		return ErrInvalidFilePath
	}
	if len(filePath) > 255 {
		return ErrInvalidFilePath
	}
	return nil
}
