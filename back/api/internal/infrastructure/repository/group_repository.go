package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/jphacks/os_2502/back/api/internal/domain/group"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/models"
)

type GroupRepositorySQLBoiler struct {
	db *sql.DB
}

func NewGroupRepositorySQLBoiler(db *sql.DB) group.Repository {
	return &GroupRepositorySQLBoiler{db: db}
}

// Groupエンティティをmodels.Groupに変換
func toGroupModel(g *group.Group) *models.Group {
	model := &models.Group{
		ID:                 g.ID(),
		OwnerUserID:        g.OwnerUserID(),
		Name:               g.Name(),
		GroupType:          string(g.GroupType()),
		Status:             string(g.Status()),
		MaxMember:          g.MaxMember(),
		CurrentMemberCount: g.CurrentMemberCount(),
		InvitationToken:    g.InvitationToken(),
		CreatedAt:          g.CreatedAt(),
		UpdatedAt:          g.UpdatedAt(),
	}

	if finalizedAt := g.FinalizedAt(); finalizedAt != nil {
		model.FinalizedAt.Valid = true
		model.FinalizedAt.Time = *finalizedAt
	}

	if countdownStartedAt := g.CountdownStartedAt(); countdownStartedAt != nil {
		model.CountdownStartedAt.Valid = true
		model.CountdownStartedAt.Time = *countdownStartedAt
	}

	if expiresAt := g.ExpiresAt(); expiresAt != nil {
		model.ExpiresAt.Valid = true
		model.ExpiresAt.Time = *expiresAt
	}

	return model
}

// models.GroupをGroupエンティティに変換
func toGroupEntity(m *models.Group) (*group.Group, error) {
	var finalizedAt *time.Time
	if m.FinalizedAt.Valid {
		t := m.FinalizedAt.Time
		finalizedAt = &t
	}

	var countdownStartedAt *time.Time
	if m.CountdownStartedAt.Valid {
		t := m.CountdownStartedAt.Time
		countdownStartedAt = &t
	}

	var expiresAt *time.Time
	if m.ExpiresAt.Valid {
		t := m.ExpiresAt.Time
		expiresAt = &t
	}

	return group.Reconstruct(
		m.ID,
		m.OwnerUserID,
		m.Name,
		group.GroupType(m.GroupType),
		group.GroupStatus(m.Status),
		m.MaxMember,
		m.CurrentMemberCount,
		m.InvitationToken,
		finalizedAt,
		countdownStartedAt,
		expiresAt,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

func (r *GroupRepositorySQLBoiler) Create(ctx context.Context, g *group.Group) error {
	model := toGroupModel(g)
	err := model.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		if db.IsDuplicateError(err) {
			return group.ErrGroupAlreadyExists
		}
		return err
	}
	return nil
}

func (r *GroupRepositorySQLBoiler) FindByID(ctx context.Context, id string) (*group.Group, error) {
	model, err := models.FindGroup(ctx, r.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, group.ErrGroupNotFound
		}
		return nil, err
	}
	return toGroupEntity(model)
}

func (r *GroupRepositorySQLBoiler) FindByInvitationToken(ctx context.Context, token string) (*group.Group, error) {
	model, err := models.Groups(
		qm.Where("invitation_token = ?", token),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, group.ErrGroupNotFound
		}
		return nil, err
	}
	return toGroupEntity(model)
}

func (r *GroupRepositorySQLBoiler) FindByOwnerUserID(ctx context.Context, ownerUserID string, limit, offset int) ([]*group.Group, error) {
	modelSlice, err := models.Groups(
		qm.Where("owner_user_id = ?", ownerUserID),
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)

	if err != nil {
		return nil, err
	}

	groups := make([]*group.Group, len(modelSlice))
	for i, m := range modelSlice {
		g, err := toGroupEntity(m)
		if err != nil {
			return nil, err
		}
		groups[i] = g
	}
	return groups, nil
}

func (r *GroupRepositorySQLBoiler) List(ctx context.Context, limit, offset int) ([]*group.Group, error) {
	modelSlice, err := models.Groups(
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)

	if err != nil {
		return nil, err
	}

	groups := make([]*group.Group, len(modelSlice))
	for i, m := range modelSlice {
		g, err := toGroupEntity(m)
		if err != nil {
			return nil, err
		}
		groups[i] = g
	}
	return groups, nil
}

func (r *GroupRepositorySQLBoiler) Update(ctx context.Context, g *group.Group) error {
	model, err := models.FindGroup(ctx, r.db, g.ID())
	if err != nil {
		if err == sql.ErrNoRows {
			return group.ErrGroupNotFound
		}
		return err
	}

	model.Name = g.Name()
	model.GroupType = string(g.GroupType())
	model.Status = string(g.Status())
	model.MaxMember = g.MaxMember()
	model.CurrentMemberCount = g.CurrentMemberCount()
	model.UpdatedAt = g.UpdatedAt()

	if finalizedAt := g.FinalizedAt(); finalizedAt != nil {
		model.FinalizedAt.Valid = true
		model.FinalizedAt.Time = *finalizedAt
	} else {
		model.FinalizedAt.Valid = false
	}

	if countdownStartedAt := g.CountdownStartedAt(); countdownStartedAt != nil {
		model.CountdownStartedAt.Valid = true
		model.CountdownStartedAt.Time = *countdownStartedAt
	} else {
		model.CountdownStartedAt.Valid = false
	}

	if expiresAt := g.ExpiresAt(); expiresAt != nil {
		model.ExpiresAt.Valid = true
		model.ExpiresAt.Time = *expiresAt
	} else {
		model.ExpiresAt.Valid = false
	}

	_, err = model.Update(ctx, r.db, boil.Infer())
	return err
}

func (r *GroupRepositorySQLBoiler) Delete(ctx context.Context, id string) error {
	model, err := models.FindGroup(ctx, r.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return group.ErrGroupNotFound
		}
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}

func (r *GroupRepositorySQLBoiler) Count(ctx context.Context) (int, error) {
	count, err := models.Groups().Count(ctx, r.db)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *GroupRepositorySQLBoiler) CountByOwnerUserID(ctx context.Context, ownerUserID string) (int, error) {
	count, err := models.Groups(
		qm.Where("owner_user_id = ?", ownerUserID),
	).Count(ctx, r.db)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *GroupRepositorySQLBoiler) FindByStatus(ctx context.Context, status string, limit, offset int) ([]*group.Group, error) {
	dbGroups, err := models.Groups(
		qm.Where("status = ?", status),
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	groups := make([]*group.Group, 0, len(dbGroups))
	for _, dbGroup := range dbGroups {
		g, err := toGroupEntity(dbGroup)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (r *GroupRepositorySQLBoiler) UpdateStatus(ctx context.Context, id string, status string) error {
	dbGroup, err := models.FindGroup(ctx, r.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return group.ErrGroupNotFound
		}
		return err
	}

	dbGroup.Status = status
	_, err = dbGroup.Update(ctx, r.db, boil.Whitelist("status", "updated_at"))
	return err
}
