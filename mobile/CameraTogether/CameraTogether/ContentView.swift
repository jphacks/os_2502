import SwiftData
import SwiftUI

struct ContentView: View {
    @Environment(\.modelContext) private var modelContext
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors
    @State private var isShowingSettings = false
    @State private var isShowingFriends = false
    @State private var isShowingNotifications = false
    @State private var hasUnreadNotifications = true

    @State private var groupManager = GroupManager()
    @State private var groupListViewModel = GroupListViewModel()

    @State private var authManager = AuthenticationManager()

    var body: some View {
        if !authManager.isAuthenticated {
            LoginView(authManager: authManager)
        } else {
            mainContent
                .task {
                    // 画面表示時にグループ一覧を取得
                    await groupManager.fetchGroups()
                }
        }
    }

    private var mainContent: some View {
        NavigationStack {
            ZStack {
                appColors.backgroundGradient
                    .ignoresSafeArea()

                ScrollView {
                    VStack(spacing: 24) {
                        Spacer()
                            .frame(height: 8)

                        NavigationLink {
                            CollageGroupMainView(authManager: authManager)
                        } label: {
                            HStack {
                                Image(systemName: "plus.circle.fill")
                                    .font(.title)
                                VStack(alignment: .leading, spacing: 2) {
                                    Text("新しいグループ")
                                        .font(.headline)
                                    Text("友達とコラージュを作成")
                                        .font(.caption)
                                        .opacity(0.8)
                                }
                                Spacer()
                                Image(systemName: "arrow.right")
                                    .font(.body)
                            }
                            .foregroundColor(appColors.textPrimary)
                            .padding(20)
                        }
                        .glassMorphism(cornerRadius: 20, opacity: 0.5)
                        .padding(.horizontal, 24)

                        if groupManager.isLoading {
                            ProgressView()
                                .padding()
                        } else if let errorMessage = groupManager.errorMessage {
                            Text(errorMessage)
                                .foregroundColor(.red)
                                .padding()
                        } else {
                            // アクティブなグループ
                            let activeGroups = groupManager.groups.filter {
                                $0.status == "recruiting" || $0.status == "ready_check"
                                    || $0.status == "countdown"
                            }

                            if !activeGroups.isEmpty {
                                VStack(alignment: .leading, spacing: 16) {
                                    HStack(spacing: 8) {
                                        ZStack {
                                            Circle()
                                                .fill(
                                                    LinearGradient(
                                                        colors: [.orange, .pink],
                                                        startPoint: .topLeading,
                                                        endPoint: .bottomTrailing
                                                    )
                                                )
                                                .frame(width: 28, height: 28)
                                            Image(systemName: "flame.fill")
                                                .font(.system(size: 14))
                                                .foregroundColor(.white)
                                        }
                                        Text("アクティブ")
                                            .font(.title3)
                                            .fontWeight(.bold)
                                            .foregroundColor(appColors.textPrimary)
                                    }
                                    .padding(.horizontal, 24)

                                    ScrollView(.horizontal, showsIndicators: false) {
                                        HStack(spacing: 16) {
                                            ForEach(groupListViewModel.getActiveGroups()) { group in
                                                NavigationLink {
                                                    groupDetailView(for: group)
                                                } label: {
                                                    GroupCardView(group: group)
                                                        .frame(width: 280)
                                                }
                                                .buttonStyle(.plain)
                                            }
                                        }
                                        .padding(.horizontal, 24)
                                    }
                                }
                            }

                            // 完了したグループ
                            let completedGroups = groupManager.groups.filter {
                                $0.status == "completed"
                            }

                            if !completedGroups.isEmpty {
                                VStack(alignment: .leading, spacing: 16) {
                                    HStack(spacing: 8) {
                                        ZStack {
                                            Circle()
                                                .fill(
                                                    LinearGradient(
                                                        colors: [.green, .mint],
                                                        startPoint: .topLeading,
                                                        endPoint: .bottomTrailing
                                                    )
                                                )
                                                .frame(width: 28, height: 28)
                                            Image(systemName: "checkmark")
                                                .font(.system(size: 14, weight: .bold))
                                                .foregroundColor(.white)
                                        }
                                        Text("完了")
                                            .font(.title3)
                                            .fontWeight(.bold)
                                            .foregroundColor(appColors.textPrimary)
                                    }
                                    .padding(.horizontal, 24)

                                    VStack(spacing: 12) {
                                        ForEach(completedGroups) { group in
                                            CompletedGroupCard(group: group)
                                                .padding(.horizontal, 24)
                                        }
                                    }
                                }
                                .padding(.top, 8)
                            }
                        }

                        Spacer(minLength: 40)
                    }
                    .padding(.vertical, 16)
                }
            }
            .navigationTitle("Collage")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    HStack(spacing: 16) {
                        Button {
                            isShowingFriends = true
                        } label: {
                            Image(systemName: "person.2.fill")
                                .foregroundColor(appColors.textPrimary)
                        }

                        Button {
                            isShowingNotifications = true
                            hasUnreadNotifications = false
                        } label: {
                            ZStack(alignment: .topTrailing) {
                                Image(systemName: "bell.fill")
                                    .foregroundColor(appColors.textPrimary)

                                if hasUnreadNotifications {
                                    Circle()
                                        .fill(Color.red)
                                        .frame(width: 8, height: 8)
                                        .offset(x: 4, y: -4)
                                }
                            }
                        }
                    }
                }
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button {
                        isShowingSettings = true
                    } label: {
                        Image(systemName: "gearshape.fill")
                            .foregroundColor(appColors.textPrimary)
                    }
                }
            }
            .toolbarBackground(.visible, for: .navigationBar)
            .toolbarBackground(Color.clear, for: .navigationBar)
            .sheet(isPresented: $isShowingSettings) {
                SettingsSheetView(authManager: authManager)
            }
            .sheet(isPresented: $isShowingNotifications) {
                NotificationListView()
            }
            .fullScreenCover(isPresented: $isShowingFriends) {
                FriendListView()
            }
        }
    }
    @ViewBuilder
    private func groupDetailView(for group: CollageGroup) -> some View {
        GroupDetailWrapperView(group: group)
    }

}

// MARK: - グループカードコンポーネント

struct GroupCard: View {
    let group: APIGroup
    @Environment(\.appColors) var appColors

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            // ステータスバッジ
            HStack {
                StatusBadge(status: group.status)
                Spacer()
                Image(systemName: "chevron.right")
                    .font(.caption)
                    .foregroundColor(appColors.textSecondary)
            }

            // グループ名
            Text(group.name)
                .font(.title3)
                .fontWeight(.bold)
                .foregroundColor(appColors.textPrimary)

            // メンバー情報
            HStack(spacing: 4) {
                Image(systemName: "person.2.fill")
                    .font(.caption)
                Text("\(group.currentMemberCount)/\(group.maxMember)")
                    .font(.caption)
            }
            .foregroundColor(appColors.textSecondary)

            // グループタイプ
            Text(groupTypeText(group.groupType))
                .font(.caption2)
                .foregroundColor(appColors.textSecondary)
                .padding(.horizontal, 8)
                .padding(.vertical, 4)
                .background(Color.secondary.opacity(0.2))
                .cornerRadius(8)
        }
        .frame(width: 200)
        .padding(16)
        .glassMorphism(cornerRadius: 16, opacity: 0.5)
    }

    private func groupTypeText(_ type: String) -> String {
        switch type {
        case "permanent":
            return "常設グループ"
        case "global_temporary":
            return "一時グループ"
        case "local_temporary":
            return "ローカル"
        default:
            return type
        }
    }
}

struct CompletedGroupCard: View {
    let group: APIGroup
    @Environment(\.appColors) var appColors

    var body: some View {
        HStack(spacing: 16) {
            // サムネイル（将来的に画像を表示）
            ZStack {
                RoundedRectangle(cornerRadius: 12)
                    .fill(
                        LinearGradient(
                            colors: [.green.opacity(0.6), .mint.opacity(0.6)],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(width: 80, height: 80)

                Image(systemName: "photo.fill")
                    .font(.title)
                    .foregroundColor(.white.opacity(0.8))
            }

            // グループ情報
            VStack(alignment: .leading, spacing: 6) {
                Text(group.name)
                    .font(.headline)
                    .foregroundColor(appColors.textPrimary)

                HStack(spacing: 4) {
                    Image(systemName: "person.2.fill")
                        .font(.caption)
                    Text("\(group.currentMemberCount)人")
                        .font(.caption)
                }
                .foregroundColor(appColors.textSecondary)

                Text(formatDate(group.createdAt))
                    .font(.caption2)
                    .foregroundColor(appColors.textSecondary)
            }

            Spacer()

            Image(systemName: "chevron.right")
                .foregroundColor(appColors.textSecondary)
        }
        .padding(16)
        .glassMorphism(cornerRadius: 16, opacity: 0.5)
    }

    private func formatDate(_ dateString: String) -> String {
        let formatter = ISO8601DateFormatter()
        if let date = formatter.date(from: dateString) {
            let displayFormatter = DateFormatter()
            displayFormatter.dateFormat = "yyyy/MM/dd"
            return displayFormatter.string(from: date)
        }
        return dateString
    }
}

struct StatusBadge: View {
    let status: String

    var body: some View {
        HStack(spacing: 4) {
            Circle()
                .fill(statusColor)
                .frame(width: 6, height: 6)
            Text(statusText)
                .font(.caption2)
                .fontWeight(.medium)
        }
        .foregroundColor(statusColor)
        .padding(.horizontal, 8)
        .padding(.vertical, 4)
        .background(statusColor.opacity(0.2))
        .cornerRadius(8)
    }

    private var statusText: String {
        switch status {
        case "recruiting":
            return "募集中"
        case "ready_check":
            return "準備確認中"
        case "countdown":
            return "カウントダウン"
        case "photo_taking":
            return "撮影中"
        case "completed":
            return "完了"
        case "expired":
            return "期限切れ"
        default:
            return status
        }
    }

    private var statusColor: Color {
        switch status {
        case "recruiting":
            return .blue
        case "ready_check":
            return .orange
        case "countdown":
            return .purple
        case "photo_taking":
            return .pink
        case "completed":
            return .green
        case "expired":
            return .gray
        default:
            return .secondary
        }
    }
}

struct GroupDetailWrapperView: View {
    let group: CollageGroup
    @State private var authManager = AuthenticationManager()
    @State private var viewModel: CollageGroupViewModel?

    var body: some View {
        Group {
            if let viewModel = viewModel {
                SimpleWaitingRoomView(viewModel: viewModel)
            } else {
                ProgressView()
            }
        }
        .onAppear {
            if viewModel == nil {
                viewModel = CollageGroupViewModel(authManager: authManager)
                viewModel?.currentGroup = group
            }
        }
    }
}

#Preview {
    ContentView()
        .modelContainer(for: Item.self, inMemory: true)
}
