import Foundation

enum GroupType: String, Codable {
    case temporaryLocal = "local_temporary"
    case temporaryGlobal = "global_temporary"
    case fixed = "permanent"

    var apiValue: String {
        self.rawValue
    }

    init?(apiValue: String) {
        switch apiValue {
        case "local_temporary":
            self = .temporaryLocal
        case "global_temporary":
            self = .temporaryGlobal
        case "permanent":
            self = .fixed
        default:
            return nil
        }
    }
}

enum GroupDuration {
    case temporary
    case fixed
}

enum GroupStatus: String, Codable {
    case recruiting = "recruiting"
    case readyCheck = "ready_check"
    case countdown = "countdown"
    case photoTaking = "photo_taking"
    case completed = "completed"
    case expired = "expired"

    var apiValue: String {
        self.rawValue
    }

    init?(apiValue: String) {
        self.init(rawValue: apiValue)
    }
}

struct CollageGroupMember: Identifiable, Codable {
    let id: String
    let name: String
    var isReady: Bool = false

    init(id: String = UUID().uuidString, name: String) {
        self.id = id
        self.name = name
    }
}

struct CollageGroup: Identifiable, Codable {
    let id: String
    let type: GroupType
    let maxMembers: Int
    var members: [CollageGroupMember]
    var status: GroupStatus
    let inviteCode: String
    var ownerId: String

    init(
        id: String = UUID().uuidString,
        type: GroupType,
        maxMembers: Int = 10,
        ownerId: String
    ) {
        self.id = id
        self.type = type
        self.maxMembers = maxMembers
        self.members = []
        self.status = .recruiting
        self.inviteCode = UUID().uuidString.prefix(8).uppercased()
        self.ownerId = ownerId
    }

    var canAddMember: Bool {
        members.count < maxMembers
    }

    var allMembersReady: Bool {
        !members.isEmpty && members.allSatisfy { $0.isReady }
    }
}

extension GroupDuration: Codable {}

// MARK: - API Conversion

extension CollageGroup {
    /// APIGroupからCollageGroupに変換
    init?(from apiGroup: APIGroup, members: [CollageGroupMember] = []) {
        guard let groupType = GroupType(apiValue: apiGroup.groupType),
            let status = GroupStatus(apiValue: apiGroup.status)
        else {
            return nil
        }

        self.id = apiGroup.id
        self.type = groupType
        self.maxMembers = apiGroup.maxMember
        self.members = members
        self.status = status
        self.inviteCode = apiGroup.invitationToken
        self.ownerId = apiGroup.ownerUserId
    }

    /// CollageGroupをAPI用のパラメータに変換
    func toCreateRequest(ownerUserId: String, name: String) -> CreateGroupRequest {
        CreateGroupRequest(
            ownerUserId: ownerUserId,
            name: name,
            groupType: type.apiValue,
            expiresAt: nil
        )
    }
}
