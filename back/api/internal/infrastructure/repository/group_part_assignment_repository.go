package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/group_part_assignment"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db/models"
)

type GroupPartAssignmentRepositorySQLBoiler struct {
	db *sql.DB
}

func NewGroupPartAssignmentRepositorySQLBoiler(db *sql.DB) group_part_assignment.Repository {
	return &GroupPartAssignmentRepositorySQLBoiler{db: db}
}

func toGroupPartAssignmentModel(gpa *group_part_assignment.GroupPartAssignment) *models.GroupPartAssignment {
	return &models.GroupPartAssignment{
		AssignmentID: gpa.AssignmentID().String(),
		GroupID:      gpa.GroupID().String(),
		UserID:       gpa.UserID().String(),
		PartID:       gpa.PartID().String(),
		CollageDay:   gpa.CollageDay().Format("2006-01-02"),
		AssignedAt:   gpa.AssignedAt(),
	}
}

func toGroupPartAssignmentEntity(m *models.GroupPartAssignment) (*group_part_assignment.GroupPartAssignment, error) {
	assignmentID, err := uuid.Parse(m.AssignmentID)
	if err != nil {
		return nil, err
	}
	groupID, err := uuid.Parse(m.GroupID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(m.UserID)
	if err != nil {
		return nil, err
	}
	partID, err := uuid.Parse(m.PartID)
	if err != nil {
		return nil, err
	}

	collageDay, err := time.Parse("2006-01-02", m.CollageDay)
	if err != nil {
		return nil, err
	}

	return group_part_assignment.Reconstruct(
		assignmentID,
		groupID,
		userID,
		partID,
		collageDay,
		m.AssignedAt,
	)
}

func (r *GroupPartAssignmentRepositorySQLBoiler) Save(ctx context.Context, assignment *group_part_assignment.GroupPartAssignment) error {
	model := toGroupPartAssignmentModel(assignment)
	return model.Upsert(ctx, r.db, true, []string{"assignment_id"}, boil.Infer(), boil.Infer())
}

func (r *GroupPartAssignmentRepositorySQLBoiler) FindByID(ctx context.Context, assignmentID uuid.UUID) (*group_part_assignment.GroupPartAssignment, error) {
	model, err := models.FindGroupPartAssignment(ctx, r.db, assignmentID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, group_part_assignment.ErrNotFound
		}
		return nil, err
	}
	return toGroupPartAssignmentEntity(model)
}

func (r *GroupPartAssignmentRepositorySQLBoiler) FindByGroupAndDay(ctx context.Context, groupID uuid.UUID, collageDay time.Time) ([]*group_part_assignment.GroupPartAssignment, error) {
	dayStr := collageDay.Format("2006-01-02")
	modelSlice, err := models.GroupPartAssignments(
		qm.Where("group_id = ? AND collage_day = ?", groupID.String(), dayStr),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	entities := make([]*group_part_assignment.GroupPartAssignment, len(modelSlice))
	for i, model := range modelSlice {
		entity, err := toGroupPartAssignmentEntity(model)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}

func (r *GroupPartAssignmentRepositorySQLBoiler) FindByUserGroupAndDay(ctx context.Context, userID, groupID uuid.UUID, collageDay time.Time) (*group_part_assignment.GroupPartAssignment, error) {
	dayStr := collageDay.Format("2006-01-02")
	model, err := models.GroupPartAssignments(
		qm.Where("user_id = ? AND group_id = ? AND collage_day = ?", userID.String(), groupID.String(), dayStr),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, group_part_assignment.ErrNotFound
		}
		return nil, err
	}
	return toGroupPartAssignmentEntity(model)
}

func (r *GroupPartAssignmentRepositorySQLBoiler) Delete(ctx context.Context, assignmentID uuid.UUID) error {
	model, err := models.FindGroupPartAssignment(ctx, r.db, assignmentID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return group_part_assignment.ErrNotFound
		}
		return err
	}
	_, err = model.Delete(ctx, r.db)
	return err
}
