import Foundation
import SwiftUI

@Observable
class GroupListViewModel {
    var groups: [CollageGroup]

    init() {
        // 初期化時に直接モックデータを設定
        var group1 = CollageGroup(
            type: .temporaryLocal,
            maxMembers: 5,
            ownerId: UUID().uuidString
        )
        group1.members = [
            CollageGroupMember(name: "あなた"),
            CollageGroupMember(name: "太郎"),
            CollageGroupMember(name: "花子"),
        ]
        group1.status = .recruiting

        var group2 = CollageGroup(
            type: .temporaryGlobal,
            maxMembers: 10,
            ownerId: UUID().uuidString
        )
        group2.members = [
            CollageGroupMember(name: "あなた"),
            CollageGroupMember(name: "次郎"),
        ]
        group2.status = .readyCheck

        var group3 = CollageGroup(
            type: .fixed,
            maxMembers: 4,
            ownerId: UUID().uuidString
        )
        group3.members = [
            CollageGroupMember(name: "あなた"),
            CollageGroupMember(name: "友達A"),
            CollageGroupMember(name: "友達B"),
            CollageGroupMember(name: "友達C"),
        ]
        group3.status = .completed

        self.groups = [group1, group2, group3]
    }

    func addGroup(_ group: CollageGroup) {
        groups.append(group)
    }

    func removeGroup(_ group: CollageGroup) {
        groups.removeAll { $0.id == group.id }
    }

    func getActiveGroups() -> [CollageGroup] {
        groups.filter { group in
            group.status != .completed
        }
    }

    func getCompletedGroups() -> [CollageGroup] {
        groups.filter { group in
            group.status == .completed
        }
    }
}
