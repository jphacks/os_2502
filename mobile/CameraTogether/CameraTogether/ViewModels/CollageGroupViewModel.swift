import Foundation
import SwiftUI

@Observable
class CollageGroupViewModel {
    var currentGroup: CollageGroup?
    var currentUserId: String
    var currentUserName: String

    init(userId: String = UUID().uuidString, userName: String = "ユーザー") {
        self.currentUserId = userId
        self.currentUserName = userName
    }

    func createGroup(type: GroupType, maxMembers: Int = 10) {
        var group = CollageGroup(type: type, maxMembers: maxMembers, ownerId: currentUserId)
        let member = CollageGroupMember(id: currentUserId, name: currentUserName)
        group.members.append(member)
        currentGroup = group
    }

    func joinGroup(inviteCode: String) -> Bool {
        guard var group = currentGroup else { return false }
        guard group.canAddMember else { return false }
        guard group.status == .recruiting else { return false }

        let member = CollageGroupMember(id: currentUserId, name: currentUserName)
        group.members.append(member)
        currentGroup = group
        return true
    }

    func addMember(name: String) -> Bool {
        guard var group = currentGroup else { return false }
        guard group.canAddMember else { return false }

        let member = CollageGroupMember(id: UUID().uuidString, name: name)
        group.members.append(member)
        currentGroup = group
        return true
    }

    func startReadyCheck() {
        guard var group = currentGroup else { return }
        guard group.ownerId == currentUserId else { return }

        group.status = .readyCheck
        currentGroup = group
    }

    func markReady() {
        guard var group = currentGroup else { return }
        guard let index = group.members.firstIndex(where: { $0.id == currentUserId }) else {
            return
        }

        group.members[index].isReady = true
        currentGroup = group

        if group.allMembersReady {
            startCountdown()
        }
    }

    func startCountdown() {
        guard var group = currentGroup else { return }
        group.status = .countdown
        currentGroup = group
    }

    func completeSession() {
        guard var group = currentGroup else { return }
        group.status = .completed
        currentGroup = group
    }

    func resetGroup() {
        currentGroup = nil
    }

    var isOwner: Bool {
        currentGroup?.ownerId == currentUserId
    }
}
