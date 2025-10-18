package user

import "errors"

var (
	// ErrInvalidFirebaseUID Firebase UIDが無効
	ErrInvalidFirebaseUID = errors.New("firebase uidが無効です")

	// ErrInvalidName ユーザー名が無効
	ErrInvalidName = errors.New("ユーザー名は1文字以上15文字以内である必要があります")

	// ErrInvalidUsername ユニークな公開IDが無効
	ErrInvalidUsername = errors.New("ユーザーIDは3〜30文字の英数字、アンダースコア、ハイフンで、英字で始まる必要があります")

	// ErrUsernameAlreadyExists ユーザーIDが既に使用されている
	ErrUsernameAlreadyExists = errors.New("このユーザーIDは既に使用されています")

	// ErrUserNotFound ユーザーが見つからない
	ErrUserNotFound = errors.New("ユーザーが見つかりません")

	// ErrUserAlreadyExists ユーザーが既に存在する
	ErrUserAlreadyExists = errors.New("ユーザーは既に存在します")
)
