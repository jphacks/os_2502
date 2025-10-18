package friend

import (
	"testing"
	"time"
)

func TestNewFriend(t *testing.T) {
	tests := []struct {
		name        string
		requesterID string
		addresseeID string
		wantErr     bool
	}{
		{
			name:        "valid friend request",
			requesterID: "user1",
			addresseeID: "user2",
			wantErr:     false,
		},
		{
			name:        "self friend request",
			requesterID: "user1",
			addresseeID: "user1",
			wantErr:     true,
		},
		{
			name:        "empty requester ID",
			requesterID: "",
			addresseeID: "user2",
			wantErr:     true,
		},
		{
			name:        "empty addressee ID",
			requesterID: "user1",
			addresseeID: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			friend, err := NewFriend(tt.requesterID, tt.addresseeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFriend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if friend == nil {
					t.Error("NewFriend() returned nil")
				}
				if friend.Status() != FriendStatusPending {
					t.Errorf("NewFriend() status = %v, want %v", friend.Status(), FriendStatusPending)
				}
			}
		})
	}
}

func TestFriend_Accept(t *testing.T) {
	friend, _ := NewFriend("user1", "user2")

	err := friend.Accept()
	if err != nil {
		t.Errorf("Accept() error = %v", err)
	}
	if friend.Status() != FriendStatusAccepted {
		t.Errorf("Accept() status = %v, want %v", friend.Status(), FriendStatusAccepted)
	}

	// Try to accept again
	err = friend.Accept()
	if err == nil {
		t.Error("Accept() should return error when status is not pending")
	}
}

func TestFriend_Reject(t *testing.T) {
	friend, _ := NewFriend("user1", "user2")

	err := friend.Reject()
	if err != nil {
		t.Errorf("Reject() error = %v", err)
	}
	if friend.Status() != FriendStatusRejected {
		t.Errorf("Reject() status = %v, want %v", friend.Status(), FriendStatusRejected)
	}

	// Try to reject again
	err = friend.Reject()
	if err == nil {
		t.Error("Reject() should return error when status is not pending")
	}
}

func TestFriend_IsExpired(t *testing.T) {
	// Create a pending friend request
	pending, _ := NewFriend("user1", "user2")
	if pending.IsExpired() {
		t.Error("IsExpired() should return false for new request")
	}

	// Create an old pending friend request (8 days old)
	oldTime := time.Now().Add(-8 * 24 * time.Hour)
	oldPending, _ := Reconstruct(
		"id",
		"user1",
		"user2",
		FriendStatusPending,
		oldTime,
		oldTime,
	)
	if !oldPending.IsExpired() {
		t.Error("IsExpired() should return true for old pending request")
	}

	// Accepted friend should not expire
	accepted, _ := NewFriend("user1", "user2")
	_ = accepted.Accept()
	if accepted.IsExpired() {
		t.Error("IsExpired() should return false for accepted friendship")
	}
}
