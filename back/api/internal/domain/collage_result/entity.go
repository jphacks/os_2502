package collage_result

import (
	"time"

	"github.com/google/uuid"
)

// CollageResult represents a collage result
type CollageResult struct {
	resultID         uuid.UUID
	templateID       uuid.UUID
	groupID          string
	fileURL          string
	targetUserNumber int
	isNotification   bool
	createdAt        time.Time
}

// NewCollageResult creates a new collage result
func NewCollageResult(templateID uuid.UUID, groupID, fileURL string, targetUserNumber int) (*CollageResult, error) {
	if templateID == uuid.Nil {
		return nil, ErrInvalidTemplateID
	}

	if groupID == "" {
		return nil, ErrInvalidGroupID
	}

	if err := validateFileURL(fileURL); err != nil {
		return nil, err
	}

	if targetUserNumber <= 0 {
		return nil, ErrInvalidTargetUserNumber
	}

	return &CollageResult{
		resultID:         uuid.New(),
		templateID:       templateID,
		groupID:          groupID,
		fileURL:          fileURL,
		targetUserNumber: targetUserNumber,
		isNotification:   false,
		createdAt:        time.Now(),
	}, nil
}

// Reconstruct reconstructs a CollageResult from repository data
func Reconstruct(
	resultID uuid.UUID,
	templateID uuid.UUID,
	groupID string,
	fileURL string,
	targetUserNumber int,
	isNotification bool,
	createdAt time.Time,
) (*CollageResult, error) {
	return &CollageResult{
		resultID:         resultID,
		templateID:       templateID,
		groupID:          groupID,
		fileURL:          fileURL,
		targetUserNumber: targetUserNumber,
		isNotification:   isNotification,
		createdAt:        createdAt,
	}, nil
}

// Getters
func (cr *CollageResult) ResultID() uuid.UUID {
	return cr.resultID
}

func (cr *CollageResult) TemplateID() uuid.UUID {
	return cr.templateID
}

func (cr *CollageResult) GroupID() string {
	return cr.groupID
}

func (cr *CollageResult) FileURL() string {
	return cr.fileURL
}

func (cr *CollageResult) TargetUserNumber() int {
	return cr.targetUserNumber
}

func (cr *CollageResult) IsNotification() bool {
	return cr.isNotification
}

func (cr *CollageResult) CreatedAt() time.Time {
	return cr.createdAt
}

// MarkAsNotified marks the result as notified
func (cr *CollageResult) MarkAsNotified() {
	cr.isNotification = true
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
