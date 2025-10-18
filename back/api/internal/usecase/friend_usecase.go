package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/friend"
	"github.com/jphacks/os_2502/back/api/internal/domain/user"
)

type FriendUseCase struct {
	friendRepo friend.Repository
	userRepo   user.Repository
}

func NewFriendUseCase(friendRepo friend.Repository, userRepo user.Repository) *FriendUseCase {
	return &FriendUseCase{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

// SendFriendRequest sends a friend request
func (uc *FriendUseCase) SendFriendRequest(ctx context.Context, requesterID, addresseeID string) (*friend.Friend, error) {
	// 申請先のユーザーが存在するかチェック
	addresseeUUID, err := uuid.Parse(addresseeID)
	if err != nil {
		return nil, user.ErrUserNotFound
	}
	if _, err := uc.userRepo.FindByID(ctx, addresseeUUID); err != nil {
		return nil, err
	}

	// 既にフレンド関係があるかチェック
	isFriend, err := uc.friendRepo.CheckFriendship(ctx, requesterID, addresseeID)
	if err != nil {
		return nil, err
	}
	if isFriend {
		return nil, friend.ErrAlreadyFriends
	}

	// 既にリクエストがあるかチェック（双方向）
	_, err = uc.friendRepo.FindByRequesterAndAddressee(ctx, requesterID, addresseeID)
	if err == nil {
		return nil, friend.ErrFriendRequestAlreadyExists
	}
	_, err = uc.friendRepo.FindByRequesterAndAddressee(ctx, addresseeID, requesterID)
	if err == nil {
		return nil, friend.ErrFriendRequestAlreadyExists
	}

	// 新しいフレンドリクエストを作成
	newFriend, err := friend.NewFriend(requesterID, addresseeID)
	if err != nil {
		return nil, err
	}

	if err := uc.friendRepo.Create(ctx, newFriend); err != nil {
		return nil, err
	}

	return newFriend, nil
}

// AcceptFriendRequest accepts a friend request
func (uc *FriendUseCase) AcceptFriendRequest(ctx context.Context, requestID, userID string) (*friend.Friend, error) {
	// リクエストを取得
	friendRequest, err := uc.friendRepo.FindByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	// 申請相手であることを確認
	if friendRequest.AddresseeID() != userID {
		return nil, friend.ErrFriendRequestNotFound
	}

	// 承認
	if err := friendRequest.Accept(); err != nil {
		return nil, err
	}

	// 更新
	if err := uc.friendRepo.Update(ctx, friendRequest); err != nil {
		return nil, err
	}

	return friendRequest, nil
}

// RejectFriendRequest rejects a friend request
func (uc *FriendUseCase) RejectFriendRequest(ctx context.Context, requestID, userID string) error {
	// リクエストを取得
	friendRequest, err := uc.friendRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}

	// 申請相手であることを確認
	if friendRequest.AddresseeID() != userID {
		return friend.ErrFriendRequestNotFound
	}

	// 拒否
	if err := friendRequest.Reject(); err != nil {
		return err
	}

	// 更新
	if err := uc.friendRepo.Update(ctx, friendRequest); err != nil {
		return err
	}

	return nil
}

// GetFriends gets all accepted friends for a user
func (uc *FriendUseCase) GetFriends(ctx context.Context, userID string, limit, offset int) ([]*friend.Friend, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.friendRepo.FindAcceptedFriends(ctx, userID, limit, offset)
}

// GetPendingReceivedRequests gets all pending received friend requests
func (uc *FriendUseCase) GetPendingReceivedRequests(ctx context.Context, userID string, limit, offset int) ([]*friend.Friend, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.friendRepo.FindPendingReceivedRequests(ctx, userID, limit, offset)
}

// GetPendingSentRequests gets all pending sent friend requests
func (uc *FriendUseCase) GetPendingSentRequests(ctx context.Context, userID string, limit, offset int) ([]*friend.Friend, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.friendRepo.FindPendingSentRequests(ctx, userID, limit, offset)
}

// RemoveFriend removes a friendship
func (uc *FriendUseCase) RemoveFriend(ctx context.Context, userID1, userID2 string) error {
	// フレンド関係を検索（双方向）
	friendRequest, err := uc.friendRepo.FindByRequesterAndAddressee(ctx, userID1, userID2)
	if err != nil {
		// 逆方向も試す
		friendRequest, err = uc.friendRepo.FindByRequesterAndAddressee(ctx, userID2, userID1)
		if err != nil {
			return friend.ErrFriendRequestNotFound
		}
	}

	// acceptedステータスでない場合はエラー
	if !friendRequest.IsAccepted() {
		return friend.ErrFriendRequestNotFound
	}

	// 削除
	return uc.friendRepo.Delete(ctx, friendRequest.ID())
}

// CancelFriendRequest cancels a sent friend request
func (uc *FriendUseCase) CancelFriendRequest(ctx context.Context, requestID, userID string) error {
	// リクエストを取得
	friendRequest, err := uc.friendRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}

	// 申請者であることを確認
	if friendRequest.RequesterID() != userID {
		return friend.ErrFriendRequestNotFound
	}

	// pendingステータスでない場合はエラー
	if !friendRequest.IsPending() {
		return friend.ErrCannotAcceptNonPending
	}

	// 削除
	return uc.friendRepo.Delete(ctx, friendRequest.ID())
}

// CheckFriendship checks if two users are friends
func (uc *FriendUseCase) CheckFriendship(ctx context.Context, userID1, userID2 string) (bool, error) {
	return uc.friendRepo.CheckFriendship(ctx, userID1, userID2)
}

// CleanupExpiredRequests cleans up expired pending requests
func (uc *FriendUseCase) CleanupExpiredRequests(ctx context.Context) (int, error) {
	return uc.friendRepo.DeleteExpiredPendingRequests(ctx)
}
