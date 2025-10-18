package group_member

import "context"

type Repository interface {
	// Create creates a new group member
	Create(ctx context.Context, member *GroupMember) error

	// FindByID finds a group member by ID
	FindByID(ctx context.Context, id string) (*GroupMember, error)

	// FindByGroupIDAndUserID finds a group member by group ID and user ID
	FindByGroupIDAndUserID(ctx context.Context, groupID, userID string) (*GroupMember, error)

	// FindByGroupID finds all members in a group
	FindByGroupID(ctx context.Context, groupID string) ([]*GroupMember, error)

	// Update updates a group member
	Update(ctx context.Context, member *GroupMember) error

	// Delete deletes a group member
	Delete(ctx context.Context, id string) error

	// DeleteByGroupIDAndUserID deletes a member by group ID and user ID
	DeleteByGroupIDAndUserID(ctx context.Context, groupID, userID string) error

	// CountByGroupID counts members in a group
	CountByGroupID(ctx context.Context, groupID string) (int, error)

	// CountReadyByGroupID counts ready members in a group
	CountReadyByGroupID(ctx context.Context, groupID string) (int, error)

	// IsOwner checks if a user is the owner of a group
	IsOwner(ctx context.Context, groupID, userID string) (bool, error)
}
