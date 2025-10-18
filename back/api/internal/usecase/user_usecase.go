package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/user"
)

type UserUseCase struct {
	repo user.Repository
}

func NewUserUseCase(repo user.Repository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, firebaseUID, name string) (*user.User, error) {
	// 既存ユーザーチェック
	existing, err := uc.repo.FindByFirebaseUID(ctx, firebaseUID)
	if err == nil && existing != nil {
		return nil, user.ErrUserAlreadyExists
	}

	u, err := user.NewUser(firebaseUID, name)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (uc *UserUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *UserUseCase) GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*user.User, error) {
	return uc.repo.FindByFirebaseUID(ctx, firebaseUID)
}

func (uc *UserUseCase) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	return uc.repo.FindByUsername(ctx, username)
}

func (uc *UserUseCase) SearchUsersByUsername(ctx context.Context, query string, limit, offset int) ([]*user.User, error) {
	return uc.repo.SearchByUsername(ctx, query, limit, offset)
}

func (uc *UserUseCase) UpdateUserName(ctx context.Context, id uuid.UUID, name string) (*user.User, error) {
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := u.UpdateName(name); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (uc *UserUseCase) SetUsername(ctx context.Context, id uuid.UUID, username string) (*user.User, error) {
	// 既存のusernameチェック
	existing, err := uc.repo.FindByUsername(ctx, username)
	if err == nil && existing != nil && existing.ID() != id {
		return nil, user.ErrUsernameAlreadyExists
	}

	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := u.SetUsername(username); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UserUseCase) ListUsers(ctx context.Context, limit, offset int) ([]*user.User, error) {
	return uc.repo.List(ctx, limit, offset)
}
