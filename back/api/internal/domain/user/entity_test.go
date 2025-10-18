package user

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		firebaseUID string
		userName    string
		wantErr     bool
	}{
		{
			name:        "valid user",
			firebaseUID: "firebase123",
			userName:    "Test User",
			wantErr:     false,
		},
		{
			name:        "empty firebase UID",
			firebaseUID: "",
			userName:    "Test User",
			wantErr:     true,
		},
		{
			name:        "empty name",
			firebaseUID: "firebase123",
			userName:    "",
			wantErr:     true,
		},
		{
			name:        "name too long",
			firebaseUID: "firebase123",
			userName:    string(make([]byte, 101)),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.firebaseUID, tt.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && user == nil {
				t.Error("NewUser() returned nil user")
			}
		})
	}
}

func TestUser_SetUsername(t *testing.T) {
	user, _ := NewUser("firebase123", "Test User")

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "valid username",
			username: "testuser123",
			wantErr:  false,
		},
		{
			name:     "username too short",
			username: "ab",
			wantErr:  true,
		},
		{
			name:     "username too long",
			username: string(make([]byte, 31)),
			wantErr:  true,
		},
		{
			name:     "username with invalid characters",
			username: "test@user",
			wantErr:  true,
		},
		{
			name:     "username starts with number",
			username: "123test",
			wantErr:  true,
		},
		{
			name:     "username with underscore",
			username: "test_user",
			wantErr:  false,
		},
		{
			name:     "username with hyphen",
			username: "test-user",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.SetUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && user.Username() == nil {
				t.Error("SetUsername() did not set username")
			}
		})
	}
}

func TestUser_UpdateName(t *testing.T) {
	user, _ := NewUser("firebase123", "Test User")

	tests := []struct {
		name    string
		newName string
		wantErr bool
	}{
		{
			name:    "valid name",
			newName: "New Name",
			wantErr: false,
		},
		{
			name:    "empty name",
			newName: "",
			wantErr: true,
		},
		{
			name:    "name too long",
			newName: string(make([]byte, 101)),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.UpdateName(tt.newName)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && user.Name() != tt.newName {
				t.Errorf("UpdateName() name = %v, want %v", user.Name(), tt.newName)
			}
		})
	}
}
