import Foundation

enum FriendStatus: String, Codable {
    case pending
    case accepted
    case rejected
}

struct Friend: Identifiable, Codable {
    let id: String
    let requesterId: String
    let addresseeId: String
    var status: FriendStatus
    let createdAt: Date
    var updatedAt: Date

    // UI表示用のプロパティ
    var name: String
    var iconName: String

    init(
        id: String = UUID().uuidString,
        requesterId: String,
        addresseeId: String,
        status: FriendStatus = .pending,
        name: String = "",
        iconName: String = "person.circle.fill"
    ) {
        self.id = id
        self.requesterId = requesterId
        self.addresseeId = addresseeId
        self.status = status
        self.createdAt = Date()
        self.updatedAt = Date()
        self.name = name
        self.iconName = iconName
    }

    /// APIFriendから変換するイニシャライザ
    init?(from apiFriend: APIFriend, currentUserId: String, userName: String = "フレンド") {
        guard let status = FriendStatus(rawValue: apiFriend.status) else {
            return nil
        }

        self.id = apiFriend.id
        self.requesterId = apiFriend.requesterId
        self.addresseeId = apiFriend.addresseeId
        self.status = status

        // 日付の変換
        let formatter = ISO8601DateFormatter()
        self.createdAt = formatter.date(from: apiFriend.createdAt) ?? Date()
        self.updatedAt = formatter.date(from: apiFriend.updatedAt) ?? Date()

        // 相手のユーザー情報（現時点では名前のみ）
        self.name = userName
        self.iconName = "person.circle.fill"
    }
}
