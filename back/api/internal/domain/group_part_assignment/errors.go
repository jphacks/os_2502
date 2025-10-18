package group_part_assignment

import "errors"

var (
	// ErrInvalidAssignmentID 割り当てIDが無効
	ErrInvalidAssignmentID = errors.New("割り当てIDが無効です")

	// ErrInvalidGroupID グループIDが無効
	ErrInvalidGroupID = errors.New("グループIDが無効です")

	// ErrInvalidUserID ユーザーIDが無効
	ErrInvalidUserID = errors.New("ユーザーIDが無効です")

	// ErrInvalidPartID パーツIDが無効
	ErrInvalidPartID = errors.New("パーツIDが無効です")

	// ErrInvalidCollageDay コラージュ日が無効
	ErrInvalidCollageDay = errors.New("コラージュ日が無効です")

	// ErrGroupPartAssignmentNotFound グループパーツ割り当てが見つからない
	ErrGroupPartAssignmentNotFound = errors.New("グループパーツ割り当てが見つかりません")

	// ErrGroupPartAssignmentAlreadyExists グループパーツ割り当てが既に存在する
	ErrGroupPartAssignmentAlreadyExists = errors.New("グループパーツ割り当ては既に存在します")

	// ErrDuplicatePartAssignment 同じグループ・日付で同じパーツが既に割り当てられている
	ErrDuplicatePartAssignment = errors.New("同じグループ・日付で同じパーツが既に割り当てられています")
)
