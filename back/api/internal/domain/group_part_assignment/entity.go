package group_part_assignment

import (
	"time"

	"github.com/google/uuid"
)

type GroupPartAssignment struct {
	assignmentID uuid.UUID
	groupID      uuid.UUID
	userID       uuid.UUID
	partID       uuid.UUID
	collageDay   time.Time
	assignedAt   time.Time
}

func NewGroupPartAssignment(groupID, userID, partID uuid.UUID, collageDay time.Time) (*GroupPartAssignment, error) {
	if groupID == uuid.Nil {
		return nil, ErrInvalidGroupID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if partID == uuid.Nil {
		return nil, ErrInvalidPartID
	}
	if collageDay.IsZero() {
		return nil, ErrInvalidCollageDay
	}

	return &GroupPartAssignment{
		assignmentID: uuid.New(),
		groupID:      groupID,
		userID:       userID,
		partID:       partID,
		collageDay:   collageDay,
		assignedAt:   time.Now(),
	}, nil
}

func Reconstruct(assignmentID, groupID, userID, partID uuid.UUID, collageDay, assignedAt time.Time) (*GroupPartAssignment, error) {
	if assignmentID == uuid.Nil {
		return nil, ErrInvalidAssignmentID
	}
	if groupID == uuid.Nil {
		return nil, ErrInvalidGroupID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if partID == uuid.Nil {
		return nil, ErrInvalidPartID
	}

	return &GroupPartAssignment{
		assignmentID: assignmentID,
		groupID:      groupID,
		userID:       userID,
		partID:       partID,
		collageDay:   collageDay,
		assignedAt:   assignedAt,
	}, nil
}

// Getters
func (gpa *GroupPartAssignment) AssignmentID() uuid.UUID {
	return gpa.assignmentID
}

func (gpa *GroupPartAssignment) GroupID() uuid.UUID {
	return gpa.groupID
}

func (gpa *GroupPartAssignment) UserID() uuid.UUID {
	return gpa.userID
}

func (gpa *GroupPartAssignment) PartID() uuid.UUID {
	return gpa.partID
}

func (gpa *GroupPartAssignment) CollageDay() time.Time {
	return gpa.collageDay
}

func (gpa *GroupPartAssignment) AssignedAt() time.Time {
	return gpa.assignedAt
}
