package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/user"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/db/models"
)

type UserRepositorySQLBoiler struct {
	db *sql.DB
}

func NewUserRepositorySQLBoiler(db *sql.DB) user.Repository {
	return &UserRepositorySQLBoiler{db: db}
}

// Userエンティティをmodels.Userに変換
func toModel(u *user.User) *models.User {
	model := &models.User{
		ID:          u.ID().String(),
		FirebaseUID: u.FirebaseUID(),
		Name:        u.Name(),
		CreatedAt:   u.CreatedAt(),
		UpdatedAt:   u.UpdatedAt(),
	}

	// usernameの設定
	if username := u.Username(); username != nil {
		model.Username.Valid = true
		model.Username.String = *username
	}

	return model
}

// models.UserをUserエンティティに変換
func toEntity(m *models.User) (*user.User, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, err
	}

	var username *string
	if m.Username.Valid {
		username = &m.Username.String
	}

	return user.Reconstruct(
		id,
		m.FirebaseUID,
		m.Name,
		username,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

func (r *UserRepositorySQLBoiler) Create(ctx context.Context, u *user.User) error {
	model := toModel(u)

	err := model.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		// 重複エラーのチェック
		if db.IsDuplicateError(err) {
			return user.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *UserRepositorySQLBoiler) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	model, err := models.FindUser(ctx, r.db, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return toEntity(model)
}

func (r *UserRepositorySQLBoiler) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*user.User, error) {
	model, err := models.Users(
		qm.Where("firebase_uid = ?", firebaseUID),
	).One(ctx, r.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return toEntity(model)
}

func (r *UserRepositorySQLBoiler) Update(ctx context.Context, u *user.User) error {
	model, err := models.FindUser(ctx, r.db, u.ID().String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.ErrUserNotFound
		}
		return err
	}

	model.Name = u.Name()
	model.UpdatedAt = u.UpdatedAt()

	// usernameの設定
	if username := u.Username(); username != nil {
		model.Username.Valid = true
		model.Username.String = *username
	} else {
		model.Username.Valid = false
	}

	_, err = model.Update(ctx, r.db, boil.Whitelist(
		models.UserColumns.Name,
		models.UserColumns.Username,
		models.UserColumns.UpdatedAt,
	))

	if err != nil {
		// username重複エラーのチェック
		if db.IsDuplicateError(err) {
			return user.ErrUsernameAlreadyExists
		}
		return err
	}

	return nil
}

func (r *UserRepositorySQLBoiler) Delete(ctx context.Context, id uuid.UUID) error {
	model, err := models.FindUser(ctx, r.db, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.ErrUserNotFound
		}
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}

func (r *UserRepositorySQLBoiler) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	modelUsers, err := models.Users(
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)

	if err != nil {
		return nil, err
	}

	users := make([]*user.User, 0, len(modelUsers))
	for _, m := range modelUsers {
		u, err := toEntity(m)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepositorySQLBoiler) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	model, err := models.Users(
		qm.Where("username = ?", username),
	).One(ctx, r.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return toEntity(model)
}

func (r *UserRepositorySQLBoiler) SearchByUsername(ctx context.Context, query string, limit, offset int) ([]*user.User, error) {
	modelUsers, err := models.Users(
		qm.Where("username LIKE ?", "%"+query+"%"),
		qm.OrderBy("username ASC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)

	if err != nil {
		return nil, err
	}

	users := make([]*user.User, 0, len(modelUsers))
	for _, m := range modelUsers {
		u, err := toEntity(m)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
