import FirebaseAuth
import Foundation
import SwiftUI

@Observable
class CollageGroupViewModel {
    var currentGroup: CollageGroup?
    var isLoading: Bool = false
    var errorMessage: String?

    private let groupAPI = GroupAPIService.shared
    private let userAPI = UserAPIService.shared
    private let authManager: AuthenticationManager
    private var memberPollingTask: Task<Void, Never>?

    var currentUserId: String {
        authManager.backendUser?.id ?? authManager.user?.uid ?? ""
    }

    var currentUserName: String {
        authManager.backendUser?.name ?? authManager.user?.displayName ?? "ユーザー"
    }

    init(authManager: AuthenticationManager) {
        self.authManager = authManager
    }

    deinit {
        stopMemberPolling()
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

            // APIレスポンスからモデルに変換してメンバー情報を設定
            await setGroupFromAPI(apiGroup)

            isLoading = false
            return true
        } catch {
            // 既に参加している場合（409 Conflict）は、グループ情報を取得して再入場
            if let apiError = error as? APIError,
                case .httpError(let statusCode, _) = apiError,
                statusCode == 409
            {
                do {
                    // 招待トークンでグループ情報を取得
                    let apiGroup = try await groupAPI.getGroupByInvitationToken(
                        invitationToken: invitationToken
                    )

                    // グループ情報を設定
                    await setGroupFromAPI(apiGroup)

                    isLoading = false
                    return true
                } catch {
                    isLoading = false
                    errorMessage = "グループ情報の取得に失敗しました: \(error.localizedDescription)"
                    return false
                }
            }

            isLoading = false
            errorMessage = "グループ参加に失敗しました: \(error.localizedDescription)"
            print("Error joining group: \(error)")
            return false
        }
    }

    /// APIGroupからグループ情報を設定
    @MainActor
    private func setGroupFromAPI(_ apiGroup: APIGroup) async {
        guard let group = CollageGroup(from: apiGroup) else { return }

        do {
            // メンバー情報を取得
            let members = try await groupAPI.getGroupMembers(groupId: group.id)

            var mutableGroup = group
            // メンバー情報をCollageGroupMemberに変換（ユーザー名も取得）
            var updatedMembers: [CollageGroupMember] = []
            for member in members {
                let userName = await fetchUserName(userId: member.userId)
                updatedMembers.append(
                    CollageGroupMember(
                        id: member.userId,
                        name: userName,
                        isReady: member.readyStatus
                    )
                )
            }
            mutableGroup.members = updatedMembers
            currentGroup = mutableGroup
        } catch {
            print("Error fetching members: \(error)")
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

    /// グループのメンバーを確定（ローカル・グローバル両対応）
    func finalizeGroup() async {
        guard let group = currentGroup else { return }

        // ローカル・グローバル両方でAPIを使用して状態を同期
        isLoading = true
        errorMessage = nil

        do {
            let apiGroup = try await groupAPI.finalizeGroupMembers(
                groupId: group.id,
                userId: currentUserId
            )

            // APIレスポンスから状態を更新
            if var updatedGroup = CollageGroup(from: apiGroup) {
                updatedGroup.members = group.members

                // オーナーを自動的に準備完了にする（ローカルでも反映）
                if let ownerIndex = updatedGroup.members.firstIndex(where: { $0.id == currentUserId }) {
                    updatedGroup.members[ownerIndex].isReady = true
                }

                currentGroup = updatedGroup
            }

            isLoading = false
        } catch {
            isLoading = false
            errorMessage = "グループ確定に失敗しました: \(error.localizedDescription)"
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
                // メンバー情報をCollageGroupMemberに変換（ユーザー名も取得）
                var updatedMembers: [CollageGroupMember] = []
                for member in members {
                    let userName = await fetchUserName(userId: member.userId)
                    updatedMembers.append(
                        CollageGroupMember(
                            id: member.userId,
                            name: userName,
                            isReady: member.readyStatus
                        )
                    )
                }
                mutableGroup.members = updatedMembers
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

    /// メンバーのポーリングを開始
    func startMemberPolling() {
        stopMemberPolling()

        memberPollingTask = Task { [weak self] in
            while !Task.isCancelled {
                guard let self = self, let groupId = self.currentGroup?.id else {
                    try? await Task.sleep(nanoseconds: 3_000_000_000)
                    continue
                }

                await self.refreshGroupMembers(groupId: groupId)

                try? await Task.sleep(nanoseconds: 3_000_000_000)
            }
        }
    }

    /// メンバーのポーリングを停止
    func stopMemberPolling() {
        memberPollingTask?.cancel()
        memberPollingTask = nil
    }

    /// グループメンバーを再取得（グループ全体の状態も更新）
    @MainActor
    private func refreshGroupMembers(groupId: String) async {
        do {
            // グループ全体の情報を取得してステータスも更新
            let apiGroup = try await groupAPI.getGroup(id: groupId)
            let members = try await groupAPI.getGroupMembers(groupId: groupId)

            // APIグループからCollageGroupに変換
            guard var updatedGroup = CollageGroup(from: apiGroup) else {
                return
            }

            // 各メンバーのユーザー情報を取得
            var updatedMembers: [CollageGroupMember] = []
            for member in members {
                let userName = await fetchUserName(userId: member.userId)
                updatedMembers.append(
                    CollageGroupMember(
                        id: member.userId,
                        name: userName,
                        isReady: member.readyStatus
                    )
                )
            }

            updatedGroup.members = updatedMembers
            currentGroup = updatedGroup

            print("🔄 Refreshed group: status=\(updatedGroup.status), finalized=\(updatedGroup.isFinalized), members=\(updatedMembers.count)")
        } catch {
            // ポーリング中のエラーは無視（UI更新しない）
            print("⚠️ Polling error (ignored): \(error.localizedDescription)")
        }
    }

    /// ユーザー名を取得（エラー時はデフォルト名を返す）
    private func fetchUserName(userId: String) async -> String {
        do {
            let user = try await userAPI.getUser(id: userId)
            return user.name
        } catch {
            return "User \(userId.prefix(8))"
        }
    }

    /// カウントダウン開始
    @MainActor
    func startCountdownWithAPI() async -> Bool {
        guard let group = currentGroup else {
            return false
        }

        // ローカルグループの場合はAPIをスキップ
        if group.type == .temporaryLocal {
            return true
        }

        isLoading = true
        errorMessage = nil

        do {
            // カウントダウン開始APIを呼び出し
            let apiGroup = try await groupAPI.startCountdown(
                groupId: group.id,
                userId: currentUserId
            )

            // グループ情報を更新
            if var updatedGroup = CollageGroup(from: apiGroup) {
                updatedGroup.members = group.members  // メンバー情報は保持
                currentGroup = updatedGroup
            }

            isLoading = false
            return true
        } catch {
            isLoading = false
            errorMessage = "撮影開始に失敗しました: \(error.localizedDescription)"
            return false
        }
    }
}
