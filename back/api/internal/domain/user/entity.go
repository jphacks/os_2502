package user

import (
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type User struct {
	id          uuid.UUID
	firebaseUID string
	name        string
	username    *string // ユニークな公開ID（オプショナル）
	createdAt   time.Time
	updatedAt   time.Time
}

func NewUser(firebaseUID, name string) (*User, error) {
	if firebaseUID == "" {
		return nil, ErrInvalidFirebaseUID
	}
	// 文字数をルーン数でカウント
	if name == "" || utf8.RuneCountInString(name) > 15 {
		return nil, ErrInvalidName
	}

	now := time.Now()
	return &User{
		id:          uuid.New(),
		firebaseUID: firebaseUID,
		name:        name,
		username:    nil, // 最初はnil、後で設定可能
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Reconstruct は既存のユーザーを復元（リポジトリから取得時に使用）
func Reconstruct(id uuid.UUID, firebaseUID, name string, username *string, createdAt, updatedAt time.Time) (*User, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidFirebaseUID
	}
	if firebaseUID == "" {
		return nil, ErrInvalidFirebaseUID
	}
	if name == "" || utf8.RuneCountInString(name) > 15 {
		return nil, ErrInvalidName
	}

	// usernameのバリデーション（設定されている場合）
	if username != nil {
		if err := validateUsername(*username); err != nil {
			return nil, err
		}
	}

	return &User{
		id:          id,
		firebaseUID: firebaseUID,
		name:        name,
		username:    username,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

// Getters
func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) FirebaseUID() string {
	return u.firebaseUID
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Username() *string {
	return u.username
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// ユーザー名を更新
func (u *User) UpdateName(name string) error {
	// 文字数をルーン数でカウント
	if name == "" || utf8.RuneCountInString(name) > 15 {
		return ErrInvalidName
	}
	u.name = name
	u.updatedAt = time.Now()
	return nil
}

// SetUsername はユニークな公開IDを設定
func (u *User) SetUsername(username string) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	u.username = &username
	u.updatedAt = time.Now()
	return nil
}

// validateUsername はusernameのバリデーション
func validateUsername(username string) error {
	if username == "" {
		return ErrInvalidUsername
	}

	// 文字数チェック（3〜30文字）
	length := utf8.RuneCountInString(username)
	if length < 3 || length > 30 {
		return ErrInvalidUsername
	}

	// 使用可能文字チェック（英数字、アンダースコア、ハイフン）
	for _, r := range username {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '-') {
			return ErrInvalidUsername
		}
	}

	// 先頭が英字であることをチェック
	firstChar := rune(username[0])
	if !((firstChar >= 'a' && firstChar <= 'z') || (firstChar >= 'A' && firstChar <= 'Z')) {
		return ErrInvalidUsername
	}

	return nil
}
