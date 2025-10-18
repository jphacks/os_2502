package friend

import "context"

type Repository interface {
	// Create creates a new friend request
	Create(ctx context.Context, friend *Friend) error

	// FindByID finds a friend request by ID
	FindByID(ctx context.Context, id string) (*Friend, error)

	// FindByRequesterAndAddressee finds a friend request by requester and addressee
	FindByRequesterAndAddressee(ctx context.Context, requesterID, addresseeID string) (*Friend, error)

	// Update updates a friend request
	Update(ctx context.Context, friend *Friend) error

	// Delete deletes a friend request
	Delete(ctx context.Context, id string) error

	// FindAcceptedFriends finds all accepted friends for a user
	FindAcceptedFriends(ctx context.Context, userID string, limit, offset int) ([]*Friend, error)

	// FindPendingReceivedRequests finds all pending received friend requests for a user
	FindPendingReceivedRequests(ctx context.Context, userID string, limit, offset int) ([]*Friend, error)

	// FindPendingSentRequests finds all pending sent friend requests for a user
	FindPendingSentRequests(ctx context.Context, userID string, limit, offset int) ([]*Friend, error)

	// CheckFriendship checks if two users are friends
	CheckFriendship(ctx context.Context, userID1, userID2 string) (bool, error)

	// DeleteExpiredPendingRequests deletes expired pending requests (older than 7 days)
	DeleteExpiredPendingRequests(ctx context.Context) (int, error)
}
