package friend

import "errors"

var (
	// ErrCannotFriendSelf cannot send friend request to self
	ErrCannotFriendSelf = errors.New("自分自身にフレンド申請はできません")

	// ErrInvalidUserID user ID is invalid
	ErrInvalidUserID = errors.New("ユーザーIDが無効です")

	// ErrCannotAcceptNonPending cannot accept non-pending request
	ErrCannotAcceptNonPending = errors.New("承認待ち以外のリクエストは承認できません")

	// ErrCannotRejectNonPending cannot reject non-pending request
	ErrCannotRejectNonPending = errors.New("承認待ち以外のリクエストは拒否できません")

	// ErrFriendRequestNotFound friend request not found
	ErrFriendRequestNotFound = errors.New("フレンドリクエストが見つかりません")

	// ErrFriendRequestAlreadyExists friend request already exists
	ErrFriendRequestAlreadyExists = errors.New("フレンドリクエストは既に存在します")

	// ErrAlreadyFriends already friends
	ErrAlreadyFriends = errors.New("既にフレンドです")
)
