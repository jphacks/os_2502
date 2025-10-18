package group_part_assignment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, assignment *GroupPartAssignment) error
	FindByID(ctx context.Context, assignmentID uuid.UUID) (*GroupPartAssignment, error)
	FindByGroupAndDay(ctx context.Context, groupID uuid.UUID, collageDay time.Time) ([]*GroupPartAssignment, error)
	FindByUserGroupAndDay(ctx context.Context, userID, groupID uuid.UUID, collageDay time.Time) (*GroupPartAssignment, error)
	Delete(ctx context.Context, assignmentID uuid.UUID) error
}
