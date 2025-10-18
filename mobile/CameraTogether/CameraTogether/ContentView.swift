import SwiftData
import SwiftUI

struct ContentView: View {
    @Environment(\.modelContext) private var modelContext
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors
    @State private var isShowingSettings = false
    @State private var groupManager = GroupManager()
    @State private var authManager = AuthenticationManager()

    var body: some View {
        if !authManager.isAuthenticated {
            // LoginView(authManager: authManager)
            mainContent
                .task {
                    // 画面表示時にグループ一覧を取得
                    await groupManager.fetchGroups()
                }
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
                            // CollageGroupMainView()
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
                                $0.status == "recruiting" || $0.status == "ready_check" || $0.status == "countdown"
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
                                            ForEach(activeGroups) { group in
                                                Text(group.name)
                                                    .padding()
                                                    .background(Color.blue.opacity(0.2))
                                                    .cornerRadius(12)
                                            }
                                        }
                                        .padding(.horizontal, 24)
                                    }
                                }
                            }

                            // 完了したグループ
                            let completedGroups = groupManager.groups.filter { $0.status == "completed" }

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
                                            Text(group.name)
                                                .padding()
                                                .background(Color.green.opacity(0.2))
                                                .cornerRadius(12)
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
        }
    }

}

#Preview {
    ContentView()
        .modelContainer(for: Item.self, inMemory: true)
}
