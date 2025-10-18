package template_part

import "errors"

var (
	ErrInvalidPartID      = errors.New("invalid part ID")
	ErrInvalidTemplateID  = errors.New("invalid template ID")
	ErrInvalidPartNumber  = errors.New("invalid part number")
	ErrInvalidDimensions  = errors.New("invalid dimensions")
	ErrNotFound           = errors.New("template part not found")
	ErrAlreadyExists      = errors.New("template part already exists")
)
