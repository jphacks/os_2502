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

// MARK: - Group Request Models

struct CreateGroupRequest: Codable {
    let ownerUserId: String
    let name: String
    let groupType: String
    let expiresAt: String?

    enum CodingKeys: String, CodingKey {
        case ownerUserId = "owner_user_id"
        case name
        case groupType = "group_type"
        case expiresAt = "expires_at"
    }
}

struct JoinGroupRequest: Codable {
    let userId: String

    enum CodingKeys: String, CodingKey {
        case userId = "user_id"
    }
}

struct FinalizeGroupRequest: Codable {
    let userId: String

    enum CodingKeys: String, CodingKey {
        case userId = "user_id"
    }
}

struct MarkReadyRequest: Codable {
    let userId: String

    enum CodingKeys: String, CodingKey {
        case userId = "user_id"
    }
}

// MARK: - Group Member Models

struct GroupMember: Codable, Identifiable {
    let id: String
    let groupId: String
    let userId: String
    let isOwner: Bool
    let readyStatus: Bool
    let readyAt: String?
    let joinedAt: String

    enum CodingKeys: String, CodingKey {
        case id
        case groupId = "group_id"
        case userId = "user_id"
        case isOwner = "is_owner"
        case readyStatus = "ready_status"
        case readyAt = "ready_at"
        case joinedAt = "joined_at"
    }
}

struct GroupMemberListResponse: Codable {
    let members: [GroupMember]
    let count: Int
}
