package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/group_part_assignment"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/models"
)

type GroupPartAssignmentRepository struct {
	db *sql.DB
}

func NewGroupPartAssignmentRepository(db *sql.DB) group_part_assignment.Repository {
	return &GroupPartAssignmentRepository{db: db}
}

// Model to Entity conversion
func toGroupPartAssignmentEntity(m *models.GroupPartAssignment) (*group_part_assignment.GroupPartAssignment, error) {
	assignmentID, err := uuid.Parse(m.AssignmentID)
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

	return group_part_assignment.Reconstruct(
		assignmentID,
		m.GroupID,
		userID,
		partID,
		m.CollageDay,
		m.AssignedAt,
	)
}

// Entity to Model conversion
func toGroupPartAssignmentModel(gpa *group_part_assignment.GroupPartAssignment) *models.GroupPartAssignment {
	return &models.GroupPartAssignment{
		AssignmentID: gpa.AssignmentID().String(),
		GroupID:      gpa.GroupID(),
		UserID:       gpa.UserID().String(),
		PartID:       gpa.PartID().String(),
		CollageDay:   gpa.CollageDay(),
		AssignedAt:   gpa.AssignedAt(),
	}
}

func (r *GroupPartAssignmentRepository) Create(ctx context.Context, assignment *group_part_assignment.GroupPartAssignment) error {
	model := toGroupPartAssignmentModel(assignment)
	return model.Insert(ctx, r.db, boil.Infer())
}

func (r *GroupPartAssignmentRepository) FindByID(ctx context.Context, assignmentID uuid.UUID) (*group_part_assignment.GroupPartAssignment, error) {
	model, err := models.FindGroupPartAssignment(ctx, r.db, assignmentID.String())
	if err == sql.ErrNoRows {
		return nil, group_part_assignment.ErrGroupPartAssignmentNotFound
	}
	if err != nil {
		return nil, err
	}
	return toGroupPartAssignmentEntity(model)
}

func (r *GroupPartAssignmentRepository) FindByGroupAndDay(ctx context.Context, groupID string, collageDay time.Time) ([]*group_part_assignment.GroupPartAssignment, error) {
	modelSlice, err := models.GroupPartAssignments(
		qm.Where("group_id = ? AND DATE(collage_day) = DATE(?)", groupID, collageDay),
		qm.OrderBy("assigned_at"),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	assignments := make([]*group_part_assignment.GroupPartAssignment, len(modelSlice))
	for i, model := range modelSlice {
		gpa, err := toGroupPartAssignmentEntity(model)
		if err != nil {
			return nil, err
		}
		assignments[i] = gpa
	}
	return assignments, nil
}

func (r *GroupPartAssignmentRepository) FindByUserGroupAndDay(ctx context.Context, userID uuid.UUID, groupID string, collageDay time.Time) (*group_part_assignment.GroupPartAssignment, error) {
	model, err := models.GroupPartAssignments(
		qm.Where("user_id = ? AND group_id = ? AND DATE(collage_day) = DATE(?)", userID.String(), groupID, collageDay),
	).One(ctx, r.db)
	if err == sql.ErrNoRows {
		return nil, group_part_assignment.ErrGroupPartAssignmentNotFound
	}
	if err != nil {
		return nil, err
	}
	return toGroupPartAssignmentEntity(model)
}

func (r *GroupPartAssignmentRepository) FindByPartID(ctx context.Context, partID uuid.UUID) ([]*group_part_assignment.GroupPartAssignment, error) {
	modelSlice, err := models.GroupPartAssignments(
		qm.Where("part_id = ?", partID.String()),
		qm.OrderBy("assigned_at"),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	assignments := make([]*group_part_assignment.GroupPartAssignment, len(modelSlice))
	for i, model := range modelSlice {
		gpa, err := toGroupPartAssignmentEntity(model)
		if err != nil {
			return nil, err
		}
		assignments[i] = gpa
	}
	return assignments, nil
}

func (r *GroupPartAssignmentRepository) Update(ctx context.Context, assignment *group_part_assignment.GroupPartAssignment) error {
	model, err := models.FindGroupPartAssignment(ctx, r.db, assignment.AssignmentID().String())
	if err == sql.ErrNoRows {
		return group_part_assignment.ErrGroupPartAssignmentNotFound
	}
	if err != nil {
		return err
	}

	model.GroupID = assignment.GroupID()
	model.UserID = assignment.UserID().String()
	model.PartID = assignment.PartID().String()
	model.CollageDay = assignment.CollageDay()

	_, err = model.Update(ctx, r.db, boil.Whitelist(
		models.GroupPartAssignmentColumns.GroupID,
		models.GroupPartAssignmentColumns.UserID,
		models.GroupPartAssignmentColumns.PartID,
		models.GroupPartAssignmentColumns.CollageDay,
	))
	return err
}

func (r *GroupPartAssignmentRepository) Delete(ctx context.Context, assignmentID uuid.UUID) error {
	model, err := models.FindGroupPartAssignment(ctx, r.db, assignmentID.String())
	if err == sql.ErrNoRows {
		return group_part_assignment.ErrGroupPartAssignmentNotFound
	}
	if err != nil {
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}

func (r *GroupPartAssignmentRepository) DeleteByGroupAndDay(ctx context.Context, groupID string, collageDay time.Time) error {
	_, err := models.GroupPartAssignments(
		qm.Where("group_id = ? AND DATE(collage_day) = DATE(?)", groupID, collageDay),
	).DeleteAll(ctx, r.db)
	return err
}

func (r *GroupPartAssignmentRepository) List(ctx context.Context, limit, offset int) ([]*group_part_assignment.GroupPartAssignment, error) {
	modelSlice, err := models.GroupPartAssignments(
		qm.OrderBy("assigned_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	assignments := make([]*group_part_assignment.GroupPartAssignment, len(modelSlice))
	for i, model := range modelSlice {
		gpa, err := toGroupPartAssignmentEntity(model)
		if err != nil {
			return nil, err
		}
		assignments[i] = gpa
	}
	return assignments, nil
}
