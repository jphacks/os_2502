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

// CreateGroup creates a new group (max_member is system fixed)
func (uc *GroupUseCase) CreateGroup(ctx context.Context, ownerUserID, name string, groupType group.GroupType, expiresAt *time.Time) (*group.Group, error) {
	// グループ作成
	newGroup, err := group.NewGroup(ownerUserID, name, groupType, expiresAt)
	if err != nil {
		return nil, err
	}

	// グループを保存
	if err := uc.groupRepo.Create(ctx, newGroup); err != nil {
		return nil, err
	}

	// オーナーをメンバーとして追加
	owner, err := group_member.NewGroupMember(newGroup.ID(), ownerUserID, true)
	if err != nil {
		return nil, err
	}
	if err := uc.memberRepo.Create(ctx, owner); err != nil {
		return nil, err
	}

	// グループのメンバー数を更新
	if err := newGroup.IncrementMemberCount(); err != nil {
		return nil, err
	}
	if err := uc.groupRepo.Update(ctx, newGroup); err != nil {
		return nil, err
	}

	return newGroup, nil
}

// JoinGroup joins a group via invitation token
func (uc *GroupUseCase) JoinGroup(ctx context.Context, invitationToken, userID string) (*group.Group, error) {
	// 招待トークンでグループを検索
	g, err := uc.groupRepo.FindByInvitationToken(ctx, invitationToken)
	if err != nil {
		return nil, err
	}

	// 参加可能かチェック
	if !g.CanJoin() {
		if g.IsExpired() {
			return nil, group.ErrGroupExpired
		}
		if g.IsFull() {
			return nil, group.ErrGroupFull
		}
		return nil, group.ErrGroupNotRecruiting
	}

	// 既に参加済みかチェック
	_, err = uc.memberRepo.FindByGroupIDAndUserID(ctx, g.ID(), userID)
	if err == nil {
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

	// グループのメンバー数を更新
	if err := g.IncrementMemberCount(); err != nil {
		return nil, err
	}
	if err := uc.groupRepo.Update(ctx, g); err != nil {
		return nil, err
	}

	return g, nil
}

// FinalizeGroupMembers finalizes the group members (owner only)
func (uc *GroupUseCase) FinalizeGroupMembers(ctx context.Context, groupID, userID string) (*group.Group, error) {
	// グループを取得
	g, err := uc.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// オーナー権限チェック
	isOwner, err := uc.memberRepo.IsOwner(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, group.ErrInvalidOwnerUserID
	}

	// メンバー確定
	if err := g.FinalizeMembers(); err != nil {
		return nil, err
	}

	// 更新
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

	// 準備完了にマーク
	if err := member.MarkReady(); err != nil {
		return err
	}

	// 更新
	if err := uc.memberRepo.Update(ctx, member); err != nil {
		return err
	}

	// 全員準備完了かチェック
	return uc.checkAllMembersReady(ctx, groupID)
}

// checkAllMembersReady checks if all members are ready and starts countdown if so
func (uc *GroupUseCase) checkAllMembersReady(ctx context.Context, groupID string) error {
	// グループを取得
	g, err := uc.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return err
	}

	// ready_check ステータスでない場合は何もしない
	if g.Status() != group.GroupStatusReadyCheck {
		return nil
	}

	// 準備完了メンバー数を取得
	readyCount, err := uc.memberRepo.CountReadyByGroupID(ctx, groupID)
	if err != nil {
		return err
	}

	// 全員準備完了の場合、カウントダウン開始
	if readyCount == g.CurrentMemberCount() {
		if err := g.StartCountdown(); err != nil {
			return err
		}
		if err := uc.groupRepo.Update(ctx, g); err != nil {
			return err
		}
	}

	return nil
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
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return uc.groupRepo.FindByOwnerUserID(ctx, ownerUserID, limit, offset)
}

// GetGroupMembers retrieves all members of a group
func (uc *GroupUseCase) GetGroupMembers(ctx context.Context, groupID string) ([]*group_member.GroupMember, error) {
	return uc.memberRepo.FindByGroupID(ctx, groupID)
}

// DeleteGroup deletes a group (owner only)
func (uc *GroupUseCase) DeleteGroup(ctx context.Context, groupID, userID string) error {
	// オーナー権限チェック
	isOwner, err := uc.memberRepo.IsOwner(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return group.ErrInvalidOwnerUserID
	}

	return uc.groupRepo.Delete(ctx, groupID)
}

// LeaveGroup allows a member to leave a group
func (uc *GroupUseCase) LeaveGroup(ctx context.Context, groupID, userID string) error {
	// グループを取得
	g, err := uc.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return err
	}

	// 募集中以外は退出不可
	if g.Status() != group.GroupStatusRecruiting {
		return group.ErrGroupNotRecruiting
	}

	// オーナーは退出不可
	isOwner, err := uc.memberRepo.IsOwner(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if isOwner {
		return group.ErrInvalidOwnerUserID
	}

	// メンバーを削除
	if err := uc.memberRepo.DeleteByGroupIDAndUserID(ctx, groupID, userID); err != nil {
		return err
	}

	// グループのメンバー数を減らす
	if err := g.DecrementMemberCount(); err != nil {
		return err
	}
	if err := uc.groupRepo.Update(ctx, g); err != nil {
		return err
	}

	return nil
}

// ListGroups retrieves all groups
func (uc *GroupUseCase) ListGroups(ctx context.Context, limit, offset int) ([]*group.Group, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return uc.groupRepo.List(ctx, limit, offset)
}
