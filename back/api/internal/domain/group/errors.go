package group

import "errors"

var (
	// Basic validation errors
	ErrInvalidGroupID     = errors.New("無効なグループIDです")
	ErrInvalidOwnerUserID = errors.New("無効なオーナーユーザーIDです")
	ErrInvalidUserID      = errors.New("無効なユーザーIDです")
	ErrInvalidName        = errors.New("グループ名は1〜15文字で入力してください")
	ErrInvalidMaxMember   = errors.New("最大メンバー数は1〜100人で設定してください")
	ErrInvalidGroupType   = errors.New("無効なグループタイプです")
	ErrInvalidGroupStatus = errors.New("無効なグループステータスです")
	ErrInvalidMemberCount = errors.New("無効なメンバー数です")

	// Business logic errors
	ErrGroupAlreadyExists       = errors.New("このグループは既に存在します")
	ErrGroupNotFound            = errors.New("グループが見つかりません")
	ErrGroupFull                = errors.New("グループが満員です")
	ErrMaxMemberLessThanCurrent = errors.New("最大メンバー数は現在のメンバー数より少なく設定できません")
	ErrNoMembers                = errors.New("メンバーがいません")

	// Status transition errors
	ErrGroupNotRecruiting  = errors.New("グループは募集中ではありません")
	ErrGroupNotReadyCheck  = errors.New("グループは準備確認中ではありません")
	ErrGroupNotCountdown   = errors.New("グループはカウントダウン中ではありません")
	ErrGroupNotPhotoTaking = errors.New("グループは撮影中ではありません")

	// Token errors
	ErrInvalidInvitationToken = errors.New("無効な招待トークンです")
	ErrGroupExpired           = errors.New("グループの有効期限が切れています")
)
