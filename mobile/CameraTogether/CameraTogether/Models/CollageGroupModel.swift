import Foundation

enum GroupType {
    case temporaryLocal
    case temporaryGlobal
    case fixed
}

enum GroupDuration {
    case temporary
    case fixed
}

enum GroupStatus {
    case recruiting
    case readyCheck
    case countdown
    case completed
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

extension GroupType: Codable {}
extension GroupDuration: Codable {}
extension GroupStatus: Codable {}
