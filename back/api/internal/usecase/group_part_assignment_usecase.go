package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/group_part_assignment"
)

type GroupPartAssignmentUseCase struct {
	repo group_part_assignment.Repository
}

func NewGroupPartAssignmentUseCase(repo group_part_assignment.Repository) *GroupPartAssignmentUseCase {
	return &GroupPartAssignmentUseCase{repo: repo}
}

func (uc *GroupPartAssignmentUseCase) CreateGroupPartAssignment(
	ctx context.Context,
	groupID string,
	userID, partID uuid.UUID,
	collageDay time.Time,
) (*group_part_assignment.GroupPartAssignment, error) {
	// 同じユーザー、グループ、日付の組み合わせが既に存在するかチェック
	existing, err := uc.repo.FindByUserGroupAndDay(ctx, userID, groupID, collageDay)
	if err == nil && existing != nil {
		return nil, group_part_assignment.ErrGroupPartAssignmentAlreadyExists
	}

	gpa, err := group_part_assignment.NewGroupPartAssignment(groupID, userID, partID, collageDay)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, gpa); err != nil {
		return nil, err
	}

	return gpa, nil
}

func (uc *GroupPartAssignmentUseCase) GetGroupPartAssignmentByID(ctx context.Context, assignmentID uuid.UUID) (*group_part_assignment.GroupPartAssignment, error) {
	return uc.repo.FindByID(ctx, assignmentID)
}

func (uc *GroupPartAssignmentUseCase) GetGroupPartAssignmentsByGroupAndDay(ctx context.Context, groupID string, collageDay time.Time) ([]*group_part_assignment.GroupPartAssignment, error) {
	return uc.repo.FindByGroupAndDay(ctx, groupID, collageDay)
}

func (uc *GroupPartAssignmentUseCase) GetGroupPartAssignmentByUserGroupAndDay(ctx context.Context, userID uuid.UUID, groupID string, collageDay time.Time) (*group_part_assignment.GroupPartAssignment, error) {
	return uc.repo.FindByUserGroupAndDay(ctx, userID, groupID, collageDay)
}

func (uc *GroupPartAssignmentUseCase) GetGroupPartAssignmentsByPartID(ctx context.Context, partID uuid.UUID) ([]*group_part_assignment.GroupPartAssignment, error) {
	return uc.repo.FindByPartID(ctx, partID)
}

func (uc *GroupPartAssignmentUseCase) DeleteGroupPartAssignment(ctx context.Context, assignmentID uuid.UUID) error {
	return uc.repo.Delete(ctx, assignmentID)
}

func (uc *GroupPartAssignmentUseCase) DeleteGroupPartAssignmentsByGroupAndDay(ctx context.Context, groupID string, collageDay time.Time) error {
	return uc.repo.DeleteByGroupAndDay(ctx, groupID, collageDay)
}

func (uc *GroupPartAssignmentUseCase) ListGroupPartAssignments(ctx context.Context, limit, offset int) ([]*group_part_assignment.GroupPartAssignment, error) {
	return uc.repo.List(ctx, limit, offset)
}
