import Foundation

// MARK: - User Models

struct User: Codable {
    let id: String
    let firebaseUid: String
    let name: String
    let username: String?
    let createdAt: String
    let updatedAt: String

    enum CodingKeys: String, CodingKey {
        case id
        case firebaseUid = "firebase_uid"
        case name
        case username
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}

struct UserListResponse: Codable {
    let users: [User]
    let limit: Int
    let offset: Int
    let count: Int
}

// MARK: - Group Models

struct APIGroup: Codable, Identifiable {
    let id: String
    let ownerUserId: String
    let name: String
    let groupType: String
    let status: String
    let maxMember: Int
    let currentMemberCount: Int
    let invitationToken: String
    let finalizedAt: String?
    let countdownStartedAt: String?
    let expiresAt: String?
    let createdAt: String
    let updatedAt: String

    enum CodingKeys: String, CodingKey {
        case id
        case ownerUserId = "owner_user_id"
        case name
        case groupType = "group_type"
        case status
        case maxMember = "max_member"
        case currentMemberCount = "current_member_count"
        case invitationToken = "invitation_token"
        case finalizedAt = "finalized_at"
        case countdownStartedAt = "countdown_started_at"
        case expiresAt = "expires_at"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}

struct GroupListResponse: Codable {
    let groups: [APIGroup]
    let totalCount: Int

    enum CodingKeys: String, CodingKey {
        case groups
        case totalCount = "total_count"
    }
}
