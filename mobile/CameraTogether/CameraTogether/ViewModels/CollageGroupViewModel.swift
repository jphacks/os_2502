import FirebaseAuth
import Foundation
import SwiftUI

@Observable
class CollageGroupViewModel {
    var currentGroup: CollageGroup?
    var isLoading: Bool = false
    var errorMessage: String?

    private let groupAPI = GroupAPIService.shared
    private let authManager: AuthenticationManager

    var currentUserId: String {
        authManager.backendUser?.id ?? authManager.user?.uid ?? ""
    }

    var currentUserName: String {
        authManager.backendUser?.name ?? authManager.user?.displayName ?? "ユーザー"
    }

    init(authManager: AuthenticationManager) {
        self.authManager = authManager
    }

    /// グループ作成
    @MainActor
    func createGroup(type: GroupType, name: String, maxMembers: Int = 10) async {
        isLoading = true
        errorMessage = nil

        do {
            // APIでグループ作成
            let apiGroup = try await groupAPI.createGroup(
                ownerUserId: currentUserId,
                name: name,
                groupType: type.apiValue
            )

            // APIレスポンスからモデルに変換
            if let group = CollageGroup(from: apiGroup) {
                var mutableGroup = group
                // 現在のユーザーをメンバーに追加（サーバー側で既に追加されている場合もある）
                let member = CollageGroupMember(id: currentUserId, name: currentUserName)
                if !mutableGroup.members.contains(where: { $0.id == currentUserId }) {
                    mutableGroup.members.append(member)
                }
                currentGroup = mutableGroup
            }

            isLoading = false
        } catch {
            isLoading = false
            errorMessage = "グループ作成に失敗しました: \(error.localizedDescription)"
            print("Error creating group: \(error)")
        }
    }

    /// グループ作成（ローカルのみ - 旧バージョン）
    func createGroupLocal(type: GroupType, maxMembers: Int = 10) {
        var group = CollageGroup(type: type, maxMembers: maxMembers, ownerId: currentUserId)
        let member = CollageGroupMember(id: currentUserId, name: currentUserName)
        group.members.append(member)
        currentGroup = group
    }

    /// グループ参加
    @MainActor
    func joinGroupWithAPI(invitationToken: String) async -> Bool {
        isLoading = true
        errorMessage = nil

        do {
            // APIでグループに参加
            let apiGroup = try await groupAPI.joinGroup(
                token: invitationToken,
                userId: currentUserId
            )

            // APIレスポンスからモデルに変換
            if let group = CollageGroup(from: apiGroup) {
                // メンバー情報を取得
                let members = try await groupAPI.getGroupMembers(groupId: group.id)

                var mutableGroup = group
                // メンバー情報をCollageGroupMemberに変換
                mutableGroup.members = members.map { member in
                    CollageGroupMember(
                        id: member.userId,
                        name: "User \(member.userId.prefix(8))"  // 実際にはユーザー名を取得する必要あり
                    )
                }
                currentGroup = mutableGroup
            }

            isLoading = false
            return true
        } catch {
            isLoading = false
            errorMessage = "グループ参加に失敗しました: \(error.localizedDescription)"
            print("Error joining group: \(error)")
            return false
        }
    }

    /// グループ参加（ローカルのみ - 旧バージョン）
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

    /// メンバー確定
    @MainActor
    func finalizeMembers() async -> Bool {
        guard let group = currentGroup else { return false }
        guard group.ownerId == currentUserId else { return false }

        isLoading = true
        errorMessage = nil

        do {
            let apiGroup = try await groupAPI.finalizeGroupMembers(
                groupId: group.id,
                userId: currentUserId
            )

            if let updatedGroup = CollageGroup(from: apiGroup) {
                var mutableGroup = updatedGroup
                mutableGroup.members = group.members
                currentGroup = mutableGroup
            }

            isLoading = false
            return true
        } catch {
            isLoading = false
            errorMessage = "メンバー確定に失敗しました: \(error.localizedDescription)"
            print("Error finalizing members: \(error)")
            return false
        }
    }

    /// メンバー確定（ローカルのみ）
    func startReadyCheck() {
        guard var group = currentGroup else { return }
        guard group.ownerId == currentUserId else { return }

        group.status = .readyCheck
        currentGroup = group
    }

    /// 準備完了
    @MainActor
    func markReadyWithAPI() async -> Bool {
        guard let group = currentGroup else { return false }

        isLoading = true
        errorMessage = nil

        do {
            try await groupAPI.markMemberReady(
                groupId: group.id,
                userId: currentUserId
            )

            // ローカルでも状態を更新
            markReady()

            isLoading = false
            return true
        } catch {
            isLoading = false
            errorMessage = "準備完了に失敗しました: \(error.localizedDescription)"
            print("Error marking ready: \(error)")
            return false
        }
    }

    /// 準備完了（ローカルのみ）
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

    /// グループ取得
    @MainActor
    func fetchGroup(groupId: String) async -> Bool {
        isLoading = true
        errorMessage = nil

        do {
            let apiGroup = try await groupAPI.getGroup(id: groupId)

            if let group = CollageGroup(from: apiGroup) {
                // メンバー情報を取得
                let members = try await groupAPI.getGroupMembers(groupId: group.id)

                var mutableGroup = group
                mutableGroup.members = members.map { member in
                    CollageGroupMember(
                        id: member.userId,
                        name: "User \(member.userId.prefix(8))"
                    )
                }
                currentGroup = mutableGroup
            }

            isLoading = false
            return true
        } catch {
            isLoading = false
            errorMessage = "グループ取得に失敗しました: \(error.localizedDescription)"
            print("Error fetching group: \(error)")
            return false
        }
    }

    /// グループ離脱
    @MainActor
    func leaveGroup() async -> Bool {
        guard let group = currentGroup else { return false }

        isLoading = true
        errorMessage = nil

        do {
            try await groupAPI.leaveGroup(
                groupId: group.id,
                userId: currentUserId
            )

            currentGroup = nil
            isLoading = false
            return true
        } catch {
            isLoading = false
            errorMessage = "グループ離脱に失敗しました: \(error.localizedDescription)"
            print("Error leaving group: \(error)")
            return false
        }
    }

    /// グループ削除
    @MainActor
    func deleteGroup() async -> Bool {
        guard let group = currentGroup else { return false }
        guard group.ownerId == currentUserId else { return false }

        isLoading = true
        errorMessage = nil

        do {
            try await groupAPI.deleteGroup(
                id: group.id,
                userId: currentUserId
            )

            currentGroup = nil
            isLoading = false
            return true
        } catch {
            isLoading = false
            errorMessage = "グループ削除に失敗しました: \(error.localizedDescription)"
            print("Error deleting group: \(error)")
            return false
        }
    }

    func resetGroup() {
        currentGroup = nil
    }

    var isOwner: Bool {
        currentGroup?.ownerId == currentUserId
    }
}
