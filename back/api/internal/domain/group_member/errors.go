package group_member

import "errors"

var (
	ErrInvalidMemberID     = errors.New("無効なメンバーIDです")
	ErrInvalidGroupID      = errors.New("無効なグループIDです")
	ErrInvalidUserID       = errors.New("無効なユーザーIDです")
	ErrMemberNotFound      = errors.New("メンバーが見つかりません")
	ErrMemberAlreadyExists = errors.New("このメンバーは既に存在します")
	ErrAlreadyReady        = errors.New("既に準備完了状態です")
	ErrNotReady            = errors.New("準備完了状態ではありません")
)
