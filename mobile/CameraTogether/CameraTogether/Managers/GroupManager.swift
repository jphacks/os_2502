import Foundation

@MainActor
@Observable
class GroupManager {
    var groups: [APIGroup] = []
    var isLoading = false
    var errorMessage: String?

    /// グループ一覧を取得
    func fetchGroups(ownerUserId: String? = nil) async {
        isLoading = true
        errorMessage = nil

        do {
            groups = try await GroupAPIService.shared.getGroups(ownerUserId: ownerUserId)
        } catch {
            errorMessage = "グループの取得に失敗しました: \(error.localizedDescription)"
            if AppConfig.enableLogging {
                print("グループ取得エラー: \(error)")
            }
        }

        isLoading = false
    }

    /// グループを作成
    func createGroup(ownerUserId: String, name: String, groupType: String = "global_temporary") async throws -> APIGroup {
        isLoading = true
        errorMessage = nil

        do {
            let group = try await GroupAPIService.shared.createGroup(
                ownerUserId: ownerUserId,
                name: name,
                groupType: groupType
            )
            // 作成したグループをリストに追加
            groups.insert(group, at: 0)
            isLoading = false
            return group
        } catch {
            errorMessage = "グループの作成に失敗しました: \(error.localizedDescription)"
            isLoading = false
            throw error
        }
    }

    /// グループを削除
    func deleteGroup(id: String, userId: String) async throws {
        isLoading = true
        errorMessage = nil

        do {
            try await GroupAPIService.shared.deleteGroup(id: id, userId: userId)
            // リストから削除
            groups.removeAll { $0.id == id }
            isLoading = false
        } catch {
            errorMessage = "グループの削除に失敗しました: \(error.localizedDescription)"
            isLoading = false
            throw error
        }
    }

    /// グループをステータスでフィルタリング
    func filterGroupsByStatus(_ status: String) -> [APIGroup] {
        groups.filter { $0.status == status }
    }

    /// 自分がオーナーのグループを取得
    func getOwnedGroups(userId: String) -> [APIGroup] {
        groups.filter { $0.ownerUserId == userId }
    }
}
