package friend

import (
	"time"

	"github.com/google/uuid"
)

// FriendStatus represents the status of a friend request
type FriendStatus string

const (
	FriendStatusPending  FriendStatus = "pending"
	FriendStatusAccepted FriendStatus = "accepted"
	FriendStatusRejected FriendStatus = "rejected"
)

// Friend represents a friendship or friend request
type Friend struct {
	id          string
	requesterID string
	addresseeID string
	status      FriendStatus
	createdAt   time.Time
	updatedAt   time.Time
}

// NewFriend creates a new friend request
func NewFriend(requesterID, addresseeID string) (*Friend, error) {
	// 自分自身へのフレンド申請はNG
	if requesterID == addresseeID {
		return nil, ErrCannotFriendSelf
	}

	if requesterID == "" || addresseeID == "" {
		return nil, ErrInvalidUserID
	}

	now := time.Now()
	return &Friend{
		id:          uuid.New().String(),
		requesterID: requesterID,
		addresseeID: addresseeID,
		status:      FriendStatusPending,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Reconstruct reconstructs a Friend from repository data
func Reconstruct(
	id string,
	requesterID string,
	addresseeID string,
	status FriendStatus,
	createdAt time.Time,
	updatedAt time.Time,
) (*Friend, error) {
	return &Friend{
		id:          id,
		requesterID: requesterID,
		addresseeID: addresseeID,
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

// Getters
func (f *Friend) ID() string {
	return f.id
}

func (f *Friend) RequesterID() string {
	return f.requesterID
}

func (f *Friend) AddresseeID() string {
	return f.addresseeID
}

func (f *Friend) Status() FriendStatus {
	return f.status
}

func (f *Friend) CreatedAt() time.Time {
	return f.createdAt
}

func (f *Friend) UpdatedAt() time.Time {
	return f.updatedAt
}

// Accept accepts a friend request
func (f *Friend) Accept() error {
	if f.status != FriendStatusPending {
		return ErrCannotAcceptNonPending
	}
	f.status = FriendStatusAccepted
	f.updatedAt = time.Now()
	return nil
}

// Reject rejects a friend request
func (f *Friend) Reject() error {
	if f.status != FriendStatusPending {
		return ErrCannotRejectNonPending
	}
	f.status = FriendStatusRejected
	f.updatedAt = time.Now()
	return nil
}

// IsAccepted returns true if the friendship is accepted
func (f *Friend) IsAccepted() bool {
	return f.status == FriendStatusAccepted
}

// IsPending returns true if the request is pending
func (f *Friend) IsPending() bool {
	return f.status == FriendStatusPending
}

// IsExpired returns true if the friend request is expired (7 days)
func (f *Friend) IsExpired() bool {
	if f.status != FriendStatusPending {
		return false
	}
	expirationDuration := 7 * 24 * time.Hour // 7日間
	return time.Now().After(f.createdAt.Add(expirationDuration))
}
