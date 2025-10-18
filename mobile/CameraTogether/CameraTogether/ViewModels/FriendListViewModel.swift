import Foundation
import SwiftData

@Observable
class FriendListViewModel {
    private var friends: [Friend] = []
    private var currentUserId: UUID = UUID()

    init() {
        loadMockData()
    }

    func getFriends() -> [Friend] {
        friends.filter { $0.status == .accepted }
    }

    func getPendingRequests() -> [Friend] {
        friends.filter { $0.status == .pending }
    }

    func acceptFriend(_ friend: Friend) {
        if let index = friends.firstIndex(where: { $0.id == friend.id }) {
            friends[index].status = .accepted
            friends[index].updatedAt = Date()
        }
    }

    func rejectFriend(_ friend: Friend) {
        if let index = friends.firstIndex(where: { $0.id == friend.id }) {
            let rejectedStatus: FriendStatus = .rejected
            friends[index].status = rejectedStatus
            friends[index].updatedAt = Date()
        }
    }

    func sendFriendRequest(to addresseeId: UUID, name: String = "新しいフレンド") {
        let pendingStatus: FriendStatus = .pending
        let newFriend = Friend(requesterId: currentUserId, addresseeId: addresseeId, status: pendingStatus, name: name, iconName: "person.circle.fill")
        friends.append(newFriend)
    }

    func addFriend(name: String, iconName: String) {
        let addresseeId = UUID()
        let acceptedStatus: FriendStatus = .accepted
        let newFriend = Friend(requesterId: currentUserId, addresseeId: addresseeId, status: acceptedStatus, name: name, iconName: iconName)
        friends.append(newFriend)
    }

    func setCurrentUserId(_ userId: UUID) {
        self.currentUserId = userId
    }

    private func loadMockData() {
        let mockUser1 = UUID()
        let mockUser2 = UUID()
        let mockUser3 = UUID()

        let acceptedStatus: FriendStatus = .accepted
        let pendingStatus: FriendStatus = .pending

        friends = [
            Friend(requesterId: currentUserId, addresseeId: mockUser1, status: acceptedStatus, name: "山田太郎", iconName: "person.circle.fill"),
            Friend(requesterId: currentUserId, addresseeId: mockUser2, status: acceptedStatus, name: "佐藤花子", iconName: "person.circle.fill"),
            Friend(requesterId: mockUser3, addresseeId: currentUserId, status: pendingStatus, name: "鈴木一郎", iconName: "person.circle.fill"),
        ]
    }
}
