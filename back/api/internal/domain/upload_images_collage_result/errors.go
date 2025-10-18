package upload_images_collage_result

import "errors"

var (
	ErrInvalidImageID    = errors.New("invalid image ID")
	ErrInvalidResultID   = errors.New("invalid result ID")
	ErrInvalidDimensions = errors.New("invalid dimensions")
	ErrNotFound          = errors.New("upload images collage result not found")
	ErrAlreadyExists     = errors.New("upload images collage result already exists")
)
