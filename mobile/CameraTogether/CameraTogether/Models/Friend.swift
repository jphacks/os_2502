import Foundation
import SwiftData

enum FriendStatus: String, Codable {
    case pending
    case accepted
    case rejected
}

@Model
class Friend {
    var id: UUID
    var requesterId: UUID
    var addresseeId: UUID
    var status: FriendStatus
    var createdAt: Date
    var updatedAt: Date

    // UI表示用の一時プロパティ（後でAPI統合時に削除予定）
    var name: String
    var iconName: String

    init(id: UUID = UUID(), requesterId: UUID, addresseeId: UUID, status: FriendStatus = FriendStatus.pending, name: String = "", iconName: String = "person.circle.fill") {
        self.id = id
        self.requesterId = requesterId
        self.addresseeId = addresseeId
        self.status = status
        self.createdAt = Date()
        self.updatedAt = Date()
        self.name = name
        self.iconName = iconName
    }
}
