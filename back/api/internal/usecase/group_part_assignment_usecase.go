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

func (uc *GroupPartAssignmentUseCase) AssignPartToUser(ctx context.Context, groupID, userID, partID uuid.UUID, collageDay time.Time) (*group_part_assignment.GroupPartAssignment, error) {
	// パーツ割り当てを作成
	assignment, err := group_part_assignment.NewGroupPartAssignment(groupID, userID, partID, collageDay)
	if err != nil {
		return nil, err
	}

	// リポジトリに保存
	if err := uc.repo.Save(ctx, assignment); err != nil {
		return nil, err
	}

	return assignment, nil
}

func (uc *GroupPartAssignmentUseCase) GetAssignmentByID(ctx context.Context, assignmentID uuid.UUID) (*group_part_assignment.GroupPartAssignment, error) {
	return uc.repo.FindByID(ctx, assignmentID)
}

func (uc *GroupPartAssignmentUseCase) GetAssignmentsByGroupAndDay(ctx context.Context, groupID uuid.UUID, collageDay time.Time) ([]*group_part_assignment.GroupPartAssignment, error) {
	return uc.repo.FindByGroupAndDay(ctx, groupID, collageDay)
}

func (uc *GroupPartAssignmentUseCase) GetUserAssignment(ctx context.Context, userID, groupID uuid.UUID, collageDay time.Time) (*group_part_assignment.GroupPartAssignment, error) {
	return uc.repo.FindByUserGroupAndDay(ctx, userID, groupID, collageDay)
}

func (uc *GroupPartAssignmentUseCase) DeleteAssignment(ctx context.Context, assignmentID uuid.UUID) error {
	return uc.repo.Delete(ctx, assignmentID)
}
