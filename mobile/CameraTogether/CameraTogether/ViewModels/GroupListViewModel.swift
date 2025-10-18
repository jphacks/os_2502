import Foundation
import SwiftUI

@Observable
class GroupListViewModel {
    var groups: [CollageGroup]
    var isLoading: Bool = false
    var errorMessage: String?

    private let groupAPI = GroupAPIService.shared

    init() {
        // 初期状態は空のリスト
        self.groups = []
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

    /// APIからグループ一覧を取得
    @MainActor
    func fetchGroups(ownerUserId: String? = nil) async {
        isLoading = true
        errorMessage = nil

        do {
            let apiGroups = try await groupAPI.getGroups(
                ownerUserId: ownerUserId,
                limit: 100,
                offset: 0
            )

            // APIGroupをCollageGroupに変換
            let convertedGroups = apiGroups.compactMap { apiGroup in
                CollageGroup(from: apiGroup)
            }

            groups = convertedGroups
            isLoading = false
        } catch {
            isLoading = false
            errorMessage = "グループ一覧の取得に失敗しました: \(error.localizedDescription)"
            print("Error fetching groups: \(error)")
        }
    }
}
