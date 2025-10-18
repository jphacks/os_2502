package device_token

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewDeviceToken(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		userID      uuid.UUID
		deviceToken string
		deviceType  DeviceType
		deviceName  *string
		wantErr     bool
	}{
		{
			name:        "valid iOS token",
			userID:      userID,
			deviceToken: "token123",
			deviceType:  DeviceTypeIOS,
			deviceName:  nil,
			wantErr:     false,
		},
		{
			name:        "valid Android token",
			userID:      userID,
			deviceToken: "token456",
			deviceType:  DeviceTypeAndroid,
			deviceName:  stringPtr("Pixel 6"),
			wantErr:     false,
		},
		{
			name:        "empty user ID",
			userID:      uuid.Nil,
			deviceToken: "token123",
			deviceType:  DeviceTypeIOS,
			deviceName:  nil,
			wantErr:     true,
		},
		{
			name:        "empty device token",
			userID:      userID,
			deviceToken: "",
			deviceType:  DeviceTypeIOS,
			deviceName:  nil,
			wantErr:     true,
		},
		{
			name:        "invalid device type",
			userID:      userID,
			deviceToken: "token123",
			deviceType:  "invalid",
			deviceName:  nil,
			wantErr:     true,
		},
		{
			name:        "token too long",
			userID:      userID,
			deviceToken: string(make([]byte, 256)),
			deviceType:  DeviceTypeIOS,
			deviceName:  nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewDeviceToken(tt.userID, tt.deviceToken, tt.deviceType, tt.deviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDeviceToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if token == nil {
					t.Error("NewDeviceToken() returned nil")
				}
				if !token.IsActive() {
					t.Error("NewDeviceToken() token should be active")
				}
			}
		})
	}
}

func TestDeviceToken_ActivateDeactivate(t *testing.T) {
	userID := uuid.New()
	token, _ := NewDeviceToken(userID, "token123", DeviceTypeIOS, nil)

	if !token.IsActive() {
		t.Error("New token should be active")
	}

	token.Deactivate()
	if token.IsActive() {
		t.Error("Token should be inactive after Deactivate()")
	}

	token.Activate()
	if !token.IsActive() {
		t.Error("Token should be active after Activate()")
	}
}

func TestDeviceToken_UpdateLastUsedAt(t *testing.T) {
	userID := uuid.New()
	token, _ := NewDeviceToken(userID, "token123", DeviceTypeIOS, nil)

	lastUsed := token.LastUsedAt()
	if lastUsed == nil {
		t.Error("LastUsedAt should not be nil")
	}

	token.UpdateLastUsedAt()
	newLastUsed := token.LastUsedAt()
	if newLastUsed == nil {
		t.Error("LastUsedAt should not be nil after update")
	}
	// Note: In a real test, you might want to add a small sleep and check that the time changed
}

func stringPtr(s string) *string {
	return &s
}
