package device_token

import (
	"time"

	"github.com/google/uuid"
)

// DeviceType represents the type of device
type DeviceType string

const (
	DeviceTypeIOS     DeviceType = "ios"
	DeviceTypeAndroid DeviceType = "android"
)

// DeviceToken represents a device token for push notifications
type DeviceToken struct {
	id          uuid.UUID
	userID      uuid.UUID
	deviceToken string
	deviceType  DeviceType
	deviceName  *string
	isActive    bool
	lastUsedAt  *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

// NewDeviceToken creates a new device token
func NewDeviceToken(userID uuid.UUID, deviceToken string, deviceType DeviceType, deviceName *string) (*DeviceToken, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	if err := validateDeviceToken(deviceToken); err != nil {
		return nil, err
	}

	if err := validateDeviceType(deviceType); err != nil {
		return nil, err
	}

	now := time.Now()
	return &DeviceToken{
		id:          uuid.New(),
		userID:      userID,
		deviceToken: deviceToken,
		deviceType:  deviceType,
		deviceName:  deviceName,
		isActive:    true,
		lastUsedAt:  &now,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Reconstruct reconstructs a DeviceToken from repository data
func Reconstruct(
	id uuid.UUID,
	userID uuid.UUID,
	deviceToken string,
	deviceType DeviceType,
	deviceName *string,
	isActive bool,
	lastUsedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) (*DeviceToken, error) {
	return &DeviceToken{
		id:          id,
		userID:      userID,
		deviceToken: deviceToken,
		deviceType:  deviceType,
		deviceName:  deviceName,
		isActive:    isActive,
		lastUsedAt:  lastUsedAt,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

// Getters
func (dt *DeviceToken) ID() uuid.UUID {
	return dt.id
}

func (dt *DeviceToken) UserID() uuid.UUID {
	return dt.userID
}

func (dt *DeviceToken) DeviceToken() string {
	return dt.deviceToken
}

func (dt *DeviceToken) DeviceType() DeviceType {
	return dt.deviceType
}

func (dt *DeviceToken) DeviceName() *string {
	return dt.deviceName
}

func (dt *DeviceToken) IsActive() bool {
	return dt.isActive
}

func (dt *DeviceToken) LastUsedAt() *time.Time {
	return dt.lastUsedAt
}

func (dt *DeviceToken) CreatedAt() time.Time {
	return dt.createdAt
}

func (dt *DeviceToken) UpdatedAt() time.Time {
	return dt.updatedAt
}

// Activate activates the device token
func (dt *DeviceToken) Activate() {
	dt.isActive = true
	dt.updatedAt = time.Now()
}

// Deactivate deactivates the device token
func (dt *DeviceToken) Deactivate() {
	dt.isActive = false
	dt.updatedAt = time.Now()
}

// UpdateLastUsedAt updates the last used timestamp
func (dt *DeviceToken) UpdateLastUsedAt() {
	now := time.Now()
	dt.lastUsedAt = &now
	dt.updatedAt = now
}

// UpdateDeviceName updates the device name
func (dt *DeviceToken) UpdateDeviceName(deviceName *string) {
	dt.deviceName = deviceName
	dt.updatedAt = time.Now()
}

// Validation functions
func validateDeviceToken(token string) error {
	if token == "" {
		return ErrInvalidDeviceToken
	}
	if len(token) > 255 {
		return ErrInvalidDeviceToken
	}
	return nil
}

func validateDeviceType(deviceType DeviceType) error {
	if deviceType != DeviceTypeIOS && deviceType != DeviceTypeAndroid {
		return ErrInvalidDeviceType
	}
	return nil
}
