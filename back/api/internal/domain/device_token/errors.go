package device_token

import "errors"

var (
	// ErrInvalidUserID user ID is invalid
	ErrInvalidUserID = errors.New("ユーザーIDが無効です")

	// ErrInvalidDeviceToken device token is invalid
	ErrInvalidDeviceToken = errors.New("デバイストークンが無効です")

	// ErrInvalidDeviceType device type is invalid
	ErrInvalidDeviceType = errors.New("デバイスタイプが無効です（ios または android を指定してください）")

	// ErrDeviceTokenNotFound device token not found
	ErrDeviceTokenNotFound = errors.New("デバイストークンが見つかりません")

	// ErrDeviceTokenAlreadyExists device token already exists
	ErrDeviceTokenAlreadyExists = errors.New("このデバイストークンは既に登録されています")
)
