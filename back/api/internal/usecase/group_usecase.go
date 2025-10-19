package usecase

import (
	"context"
	"time"

	"github.com/jphacks/os_2502/back/api/internal/domain/group"
	"github.com/jphacks/os_2502/back/api/internal/domain/group_member"
)

type GroupUseCase struct {
	groupRepo  group.Repository
	memberRepo group_member.Repository
}

func NewGroupUseCase(groupRepo group.Repository, memberRepo group_member.Repository) *GroupUseCase {
	return &GroupUseCase{
		groupRepo:  groupRepo,
		memberRepo: memberRepo,
	}
}

// CreateGroup creates a new group and adds the owner as the first member
func (uc *GroupUseCase) CreateGroup(ctx context.Context, ownerUserID, name string, groupType group.GroupType, expiresAt *time.Time) (*group.Group, error) {
	// グループの作成
	g, err := group.NewGroup(ownerUserID, name, groupType, expiresAt)
	if err != nil {
		return nil, err
	}

	// グループをリポジトリに保存
	if err := uc.groupRepo.Create(ctx, g); err != nil {
		return nil, err
	}

	// オーナーをメンバーとして追加
	member, err := group_member.NewGroupMember(g.ID(), ownerUserID, true)
	if err != nil {
		return nil, err
	}

	if err := uc.memberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	// メンバー数を更新
	if err := g.IncrementMemberCount(); err != nil {
		return nil, err
	}

	if err := uc.groupRepo.Update(ctx, g); err != nil {
		return nil, err
	}

	return g, nil
}

// GetGroupByID retrieves a group by ID
func (uc *GroupUseCase) GetGroupByID(ctx context.Context, id string) (*group.Group, error) {
	return uc.groupRepo.FindByID(ctx, id)
}

// GetGroupByInvitationToken retrieves a group by invitation token
func (uc *GroupUseCase) GetGroupByInvitationToken(ctx context.Context, token string) (*group.Group, error) {
	return uc.groupRepo.FindByInvitationToken(ctx, token)
}

// GetGroupsByOwnerUserID retrieves groups by owner user ID
func (uc *GroupUseCase) GetGroupsByOwnerUserID(ctx context.Context, ownerUserID string, limit, offset int) ([]*group.Group, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.groupRepo.FindByOwnerUserID(ctx, ownerUserID, limit, offset)
}

// JoinGroup joins a group via invitation token
func (uc *GroupUseCase) JoinGroup(ctx context.Context, invitationToken, userID string) (*group.Group, error) {
	// グループを取得
	g, err := uc.groupRepo.FindByInvitationToken(ctx, invitationToken)
	if err != nil {
		return nil, err
	}

	// グループに参加可能かチェック
	if !g.CanJoin() {
		if g.IsExpired() {
			return nil, group.ErrGroupExpired
		}
		if g.IsFull() {
			return nil, group.ErrGroupFull
		}
		return nil, group.ErrGroupNotRecruiting
	}

	// 既にメンバーかどうかチェック
	existing, err := uc.memberRepo.FindByGroupIDAndUserID(ctx, g.ID(), userID)
	if err == nil && existing != nil {
		return nil, group_member.ErrMemberAlreadyExists
	}

	// メンバーとして追加
	member, err := group_member.NewGroupMember(g.ID(), userID, false)
	if err != nil {
		return nil, err
	}

	if err := uc.memberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	// メンバー数を更新
	if err := g.IncrementMemberCount(); err != nil {
		return nil, err
	}

	if err := uc.groupRepo.Update(ctx, g); err != nil {
		return nil, err
	}

	return g, nil
}

// FinalizeGroupMembers finalizes group members (owner only)
func (uc *GroupUseCase) FinalizeGroupMembers(ctx context.Context, groupID, userID string) (*group.Group, error) {
	// グループを取得
	g, err := uc.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// オーナーかどうかチェック
	if g.OwnerUserID() != userID {
		return nil, group.ErrInvalidOwnerUserID
	}

	// メンバーを確定
	if err := g.FinalizeMembers(); err != nil {
		return nil, err
	}

	if err := uc.groupRepo.Update(ctx, g); err != nil {
		return nil, err
	}

	return g, nil
}

// MarkMemberReady marks a member as ready
func (uc *GroupUseCase) MarkMemberReady(ctx context.Context, groupID, userID string) error {
	// メンバーを取得
	member, err := uc.memberRepo.FindByGroupIDAndUserID(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if member == nil {
		return group_member.ErrMemberNotFound
	}

	// 準備完了にする
	if err := member.MarkReady(); err != nil {
		return err
	}

	if err := uc.memberRepo.Update(ctx, member); err != nil {
		return err
	}

	// 全員準備完了かチェック
	readyCount, err := uc.memberRepo.CountReadyByGroupID(ctx, groupID)
	if err != nil {
		return err
	}

	g, err := uc.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return err
	}

	// 全員準備完了ならカウントダウン開始
	if readyCount == g.CurrentMemberCount() && g.Status() == group.GroupStatusReadyCheck {
		if err := g.StartCountdown(); err != nil {
			return err
		}
		if err := uc.groupRepo.Update(ctx, g); err != nil {
			return err
		}
	}

	return nil
}

// GetGroupMembers retrieves all members of a group
func (uc *GroupUseCase) GetGroupMembers(ctx context.Context, groupID string) ([]*group_member.GroupMember, error) {
	return uc.memberRepo.FindByGroupID(ctx, groupID)
}

// StartCountdown starts the countdown for photo session
func (uc *GroupUseCase) StartCountdown(ctx context.Context, groupID, userID string) (*group.Group, error) {
	g, err := uc.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// オーナーチェック
	if g.OwnerUserID() != userID {
		return nil, group.ErrInvalidOwnerUserID
	}

	// カウントダウン開始
	if err := g.StartCountdown(); err != nil {
		return nil, err
	}

	// 更新
	if err := uc.groupRepo.Update(ctx, g); err != nil {
		return nil, err
	}

	return g, nil
}

// LeaveGroup allows a member to leave a group
func (uc *GroupUseCase) LeaveGroup(ctx context.Context, groupID, userID string) error {
	// グループを取得
	g, err := uc.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return err
	}

	// オーナーは離脱できない
	if g.OwnerUserID() == userID {
		return group.ErrInvalidOwnerUserID
	}

	// メンバー募集中のみ離脱可能
	if g.Status() != group.GroupStatusRecruiting {
		return group.ErrGroupNotRecruiting
	}

	// メンバーを削除
	member, err := uc.memberRepo.FindByGroupIDAndUserID(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if member == nil {
		return group_member.ErrMemberNotFound
	}

	if err := uc.memberRepo.DeleteByGroupIDAndUserID(ctx, groupID, userID); err != nil {
		return err
	}

	// メンバー数を更新
	if err := g.DecrementMemberCount(); err != nil {
		return err
	}

	if err := uc.groupRepo.Update(ctx, g); err != nil {
		return err
	}

	return nil
}

// DeleteGroup deletes a group (owner only)
func (uc *GroupUseCase) DeleteGroup(ctx context.Context, groupID, userID string) error {
	// グループを取得
	g, err := uc.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return err
	}

	// オーナーかどうかチェック
	if g.OwnerUserID() != userID {
		return group.ErrInvalidOwnerUserID
	}

	// グループを削除
	return uc.groupRepo.Delete(ctx, groupID)
}

// ListGroups retrieves all groups
func (uc *GroupUseCase) ListGroups(ctx context.Context, limit, offset int) ([]*group.Group, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.groupRepo.List(ctx, limit, offset)
}
