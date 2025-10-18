package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/jphacks/os_2502/back/api/internal/domain/group_member"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/models"
)

type GroupMemberRepositorySQLBoiler struct {
	db *sql.DB
}

func NewGroupMemberRepositorySQLBoiler(db *sql.DB) group_member.Repository {
	return &GroupMemberRepositorySQLBoiler{db: db}
}

// Model to Entity conversion
func toGroupMemberEntity(m *models.GroupMember) (*group_member.GroupMember, error) {
	var readyAt *time.Time
	if m.ReadyAt.Valid {
		t := m.ReadyAt.Time
		readyAt = &t
	}

	return group_member.Reconstruct(
		m.ID,
		m.GroupID,
		m.UserID,
		m.IsOwner,
		m.ReadyStatus,
		readyAt,
		m.JoinedAt,
		m.UpdatedAt,
	)
}

// Entity to Model conversion
func toGroupMemberModel(gm *group_member.GroupMember) *models.GroupMember {
	model := &models.GroupMember{
		ID:          gm.ID(),
		GroupID:     gm.GroupID(),
		UserID:      gm.UserID(),
		IsOwner:     gm.IsOwner(),
		ReadyStatus: gm.ReadyStatus(),
		JoinedAt:    gm.JoinedAt(),
		UpdatedAt:   gm.UpdatedAt(),
	}

	if readyAt := gm.ReadyAt(); readyAt != nil {
		model.ReadyAt.Valid = true
		model.ReadyAt.Time = *readyAt
	}

	return model
}

func (r *GroupMemberRepositorySQLBoiler) Create(ctx context.Context, member *group_member.GroupMember) error {
	model := toGroupMemberModel(member)
	err := model.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		if db.IsDuplicateError(err) {
			return group_member.ErrMemberAlreadyExists
		}
		return err
	}
	return nil
}

func (r *GroupMemberRepositorySQLBoiler) FindByID(ctx context.Context, id string) (*group_member.GroupMember, error) {
	model, err := models.GroupMembers(
		qm.Where("id = ?", id),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, group_member.ErrMemberNotFound
		}
		return nil, err
	}
	return toGroupMemberEntity(model)
}

func (r *GroupMemberRepositorySQLBoiler) FindByGroupIDAndUserID(ctx context.Context, groupID, userID string) (*group_member.GroupMember, error) {
	model, err := models.GroupMembers(
		qm.Where("group_id = ? AND user_id = ?", groupID, userID),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, group_member.ErrMemberNotFound
		}
		return nil, err
	}
	return toGroupMemberEntity(model)
}

func (r *GroupMemberRepositorySQLBoiler) FindByGroupID(ctx context.Context, groupID string) ([]*group_member.GroupMember, error) {
	modelSlice, err := models.GroupMembers(
		qm.Where("group_id = ?", groupID),
		qm.OrderBy("joined_at ASC"),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	members := make([]*group_member.GroupMember, len(modelSlice))
	for i, model := range modelSlice {
		member, err := toGroupMemberEntity(model)
		if err != nil {
			return nil, err
		}
		members[i] = member
	}
	return members, nil
}

func (r *GroupMemberRepositorySQLBoiler) Update(ctx context.Context, member *group_member.GroupMember) error {
	model, err := models.GroupMembers(
		qm.Where("id = ?", member.ID()),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return group_member.ErrMemberNotFound
		}
		return err
	}

	model.ReadyStatus = member.ReadyStatus()
	model.UpdatedAt = member.UpdatedAt()

	if readyAt := member.ReadyAt(); readyAt != nil {
		model.ReadyAt.Valid = true
		model.ReadyAt.Time = *readyAt
	} else {
		model.ReadyAt.Valid = false
	}

	_, err = model.Update(ctx, r.db, boil.Infer())
	return err
}

func (r *GroupMemberRepositorySQLBoiler) Delete(ctx context.Context, id string) error {
	model, err := models.GroupMembers(
		qm.Where("id = ?", id),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return group_member.ErrMemberNotFound
		}
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}

func (r *GroupMemberRepositorySQLBoiler) DeleteByGroupIDAndUserID(ctx context.Context, groupID, userID string) error {
	model, err := models.GroupMembers(
		qm.Where("group_id = ? AND user_id = ?", groupID, userID),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return group_member.ErrMemberNotFound
		}
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}

func (r *GroupMemberRepositorySQLBoiler) CountByGroupID(ctx context.Context, groupID string) (int, error) {
	count, err := models.GroupMembers(
		qm.Where("group_id = ?", groupID),
	).Count(ctx, r.db)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *GroupMemberRepositorySQLBoiler) CountReadyByGroupID(ctx context.Context, groupID string) (int, error) {
	count, err := models.GroupMembers(
		qm.Where("group_id = ? AND ready_status = ?", groupID, true),
	).Count(ctx, r.db)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *GroupMemberRepositorySQLBoiler) IsOwner(ctx context.Context, groupID, userID string) (bool, error) {
	count, err := models.GroupMembers(
		qm.Where("group_id = ? AND user_id = ? AND is_owner = ?", groupID, userID, true),
	).Count(ctx, r.db)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
