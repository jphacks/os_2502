package group_member

import (
	"time"

	"github.com/google/uuid"
)

type GroupMember struct {
	id          string
	groupID     string
	userID      string
	isOwner     bool
	readyStatus bool
	readyAt     *time.Time
	joinedAt    time.Time
	updatedAt   time.Time
}

// NewGroupMember creates a new group member
func NewGroupMember(groupID, userID string, isOwner bool) (*GroupMember, error) {
	if groupID == "" {
		return nil, ErrInvalidGroupID
	}
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	now := time.Now()
	return &GroupMember{
		id:          uuid.New().String(),
		groupID:     groupID,
		userID:      userID,
		isOwner:     isOwner,
		readyStatus: false,
		readyAt:     nil,
		joinedAt:    now,
		updatedAt:   now,
	}, nil
}

// Reconstruct reconstructs a group member from repository
func Reconstruct(id, groupID, userID string, isOwner, readyStatus bool, readyAt *time.Time, joinedAt, updatedAt time.Time) (*GroupMember, error) {
	if id == "" {
		return nil, ErrInvalidMemberID
	}
	if groupID == "" {
		return nil, ErrInvalidGroupID
	}
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	return &GroupMember{
		id:          id,
		groupID:     groupID,
		userID:      userID,
		isOwner:     isOwner,
		readyStatus: readyStatus,
		readyAt:     readyAt,
		joinedAt:    joinedAt,
		updatedAt:   updatedAt,
	}, nil
}

// Getters
func (gm *GroupMember) ID() string {
	return gm.id
}

func (gm *GroupMember) GroupID() string {
	return gm.groupID
}

func (gm *GroupMember) UserID() string {
	return gm.userID
}

func (gm *GroupMember) IsOwner() bool {
	return gm.isOwner
}

func (gm *GroupMember) ReadyStatus() bool {
	return gm.readyStatus
}

func (gm *GroupMember) ReadyAt() *time.Time {
	return gm.readyAt
}

func (gm *GroupMember) JoinedAt() time.Time {
	return gm.joinedAt
}

func (gm *GroupMember) UpdatedAt() time.Time {
	return gm.updatedAt
}

// MarkReady marks the member as ready
func (gm *GroupMember) MarkReady() error {
	if gm.readyStatus {
		return ErrAlreadyReady
	}
	now := time.Now()
	gm.readyStatus = true
	gm.readyAt = &now
	gm.updatedAt = now
	return nil
}

// CancelReady cancels the ready status
func (gm *GroupMember) CancelReady() error {
	if !gm.readyStatus {
		return ErrNotReady
	}
	gm.readyStatus = false
	gm.readyAt = nil
	gm.updatedAt = time.Now()
	return nil
}
