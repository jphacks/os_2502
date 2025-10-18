package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/jphacks/os_2502/back/api/internal/domain/friend"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db/models"
)

type FriendRepositorySQLBoiler struct {
	db *sql.DB
}

func NewFriendRepositorySQLBoiler(db *sql.DB) friend.Repository {
	return &FriendRepositorySQLBoiler{db: db}
}

// Model to Entity conversion
func toFriendEntity(m *models.Friend) (*friend.Friend, error) {
	return friend.Reconstruct(
		m.ID,
		m.RequesterID,
		m.AddresseeID,
		friend.FriendStatus(m.Status),
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// Entity to Model conversion
func toFriendModel(f *friend.Friend) *models.Friend {
	return &models.Friend{
		ID:          f.ID(),
		RequesterID: f.RequesterID(),
		AddresseeID: f.AddresseeID(),
		Status:      string(f.Status()),
		CreatedAt:   f.CreatedAt(),
		UpdatedAt:   f.UpdatedAt(),
	}
}

func (r *FriendRepositorySQLBoiler) Create(ctx context.Context, f *friend.Friend) error {
	model := toFriendModel(f)
	err := model.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		if db.IsDuplicateError(err) {
			return friend.ErrFriendRequestAlreadyExists
		}
		return err
	}
	return nil
}

func (r *FriendRepositorySQLBoiler) FindByID(ctx context.Context, id string) (*friend.Friend, error) {
	model, err := models.Friends(
		qm.Where("id = ?", id),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, friend.ErrFriendRequestNotFound
		}
		return nil, err
	}
	return toFriendEntity(model)
}

func (r *FriendRepositorySQLBoiler) FindByRequesterAndAddressee(ctx context.Context, requesterID, addresseeID string) (*friend.Friend, error) {
	model, err := models.Friends(
		qm.Where("requester_id = ? AND addressee_id = ?", requesterID, addresseeID),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, friend.ErrFriendRequestNotFound
		}
		return nil, err
	}
	return toFriendEntity(model)
}

func (r *FriendRepositorySQLBoiler) Update(ctx context.Context, f *friend.Friend) error {
	model, err := models.Friends(
		qm.Where("id = ?", f.ID()),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return friend.ErrFriendRequestNotFound
		}
		return err
	}

	model.Status = string(f.Status())
	model.UpdatedAt = f.UpdatedAt()

	_, err = model.Update(ctx, r.db, boil.Whitelist(
		models.FriendColumns.Status,
		models.FriendColumns.UpdatedAt,
	))
	return err
}

func (r *FriendRepositorySQLBoiler) Delete(ctx context.Context, id string) error {
	model, err := models.Friends(
		qm.Where("id = ?", id),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return friend.ErrFriendRequestNotFound
		}
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}

func (r *FriendRepositorySQLBoiler) FindAcceptedFriends(ctx context.Context, userID string, limit, offset int) ([]*friend.Friend, error) {
	modelSlice, err := models.Friends(
		qm.Where("(requester_id = ? OR addressee_id = ?) AND status = ?", userID, userID, string(friend.FriendStatusAccepted)),
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	friends := make([]*friend.Friend, len(modelSlice))
	for i, model := range modelSlice {
		f, err := toFriendEntity(model)
		if err != nil {
			return nil, err
		}
		friends[i] = f
	}
	return friends, nil
}

func (r *FriendRepositorySQLBoiler) FindPendingReceivedRequests(ctx context.Context, userID string, limit, offset int) ([]*friend.Friend, error) {
	modelSlice, err := models.Friends(
		qm.Where("addressee_id = ? AND status = ?", userID, string(friend.FriendStatusPending)),
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	friends := make([]*friend.Friend, len(modelSlice))
	for i, model := range modelSlice {
		f, err := toFriendEntity(model)
		if err != nil {
			return nil, err
		}
		friends[i] = f
	}
	return friends, nil
}

func (r *FriendRepositorySQLBoiler) FindPendingSentRequests(ctx context.Context, userID string, limit, offset int) ([]*friend.Friend, error) {
	modelSlice, err := models.Friends(
		qm.Where("requester_id = ? AND status = ?", userID, string(friend.FriendStatusPending)),
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	friends := make([]*friend.Friend, len(modelSlice))
	for i, model := range modelSlice {
		f, err := toFriendEntity(model)
		if err != nil {
			return nil, err
		}
		friends[i] = f
	}
	return friends, nil
}

func (r *FriendRepositorySQLBoiler) CheckFriendship(ctx context.Context, userID1, userID2 string) (bool, error) {
	count, err := models.Friends(
		qm.Where("((requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?)) AND status = ?",
			userID1, userID2, userID2, userID1, string(friend.FriendStatusAccepted)),
	).Count(ctx, r.db)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FriendRepositorySQLBoiler) DeleteExpiredPendingRequests(ctx context.Context) (int, error) {
	expirationDate := time.Now().Add(-7 * 24 * time.Hour) // 7日前

	count, err := models.Friends(
		qm.Where("status = ? AND created_at < ?", string(friend.FriendStatusPending), expirationDate),
	).DeleteAll(ctx, r.db)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
