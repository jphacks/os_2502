package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/null/v8"
	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/user"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &UserRepository{db: db}
}

// Helper functions for null.String conversion
func nullStringFromPtr(s *string) null.String {
	if s == nil {
		return null.String{Valid: false}
	}
	return null.String{String: *s, Valid: true}
}

func ptrFromNullString(ns null.String) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// Model to Entity conversion
func toUserEntity(m *models.User) (*user.User, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, err
	}

	return user.Reconstruct(
		id,
		m.FirebaseUID,
		m.Name,
		ptrFromNullString(m.Username),
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// Entity to Model conversion
func toUserModel(u *user.User) *models.User {
	return &models.User{
		ID:          u.ID().String(),
		FirebaseUID: u.FirebaseUID(),
		Name:        u.Name(),
		Username:    nullStringFromPtr(u.Username()),
		CreatedAt:   u.CreatedAt(),
		UpdatedAt:   u.UpdatedAt(),
	}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	model := toUserModel(u)
	return model.Insert(ctx, r.db, boil.Infer())
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	model, err := models.FindUser(ctx, r.db, id.String())
	if err == sql.ErrNoRows {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return toUserEntity(model)
}

func (r *UserRepository) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*user.User, error) {
	model, err := models.Users(
		qm.Where("firebase_uid = ?", firebaseUID),
	).One(ctx, r.db)
	if err == sql.ErrNoRows {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return toUserEntity(model)
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	model, err := models.Users(
		qm.Where("username = ?", username),
	).One(ctx, r.db)
	if err == sql.ErrNoRows {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return toUserEntity(model)
}

func (r *UserRepository) SearchByUsername(ctx context.Context, query string, limit, offset int) ([]*user.User, error) {
	modelSlice, err := models.Users(
		qm.Where("username LIKE ?", "%"+query+"%"),
		qm.OrderBy("username"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	users := make([]*user.User, len(modelSlice))
	for i, model := range modelSlice {
		u, err := toUserEntity(model)
		if err != nil {
			return nil, err
		}
		users[i] = u
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	model, err := models.FindUser(ctx, r.db, u.ID().String())
	if err == sql.ErrNoRows {
		return user.ErrUserNotFound
	}
	if err != nil {
		return err
	}

	model.Name = u.Name()
	model.Username = nullStringFromPtr(u.Username())
	model.UpdatedAt = u.UpdatedAt()

	_, err = model.Update(ctx, r.db, boil.Whitelist(
		models.UserColumns.Name,
		models.UserColumns.Username,
		models.UserColumns.UpdatedAt,
	))
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	model, err := models.FindUser(ctx, r.db, id.String())
	if err == sql.ErrNoRows {
		return user.ErrUserNotFound
	}
	if err != nil {
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	modelSlice, err := models.Users(
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	users := make([]*user.User, len(modelSlice))
	for i, model := range modelSlice {
		u, err := toUserEntity(model)
		if err != nil {
			return nil, err
		}
		users[i] = u
	}
	return users, nil
}
