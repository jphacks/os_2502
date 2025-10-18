package group_part_assignment

import "errors"

var (
	ErrInvalidAssignmentID = errors.New("invalid assignment ID")
	ErrInvalidGroupID      = errors.New("invalid group ID")
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrInvalidPartID       = errors.New("invalid part ID")
	ErrInvalidCollageDay   = errors.New("invalid collage day")
	ErrNotFound            = errors.New("group part assignment not found")
	ErrAlreadyExists       = errors.New("group part assignment already exists")
)
