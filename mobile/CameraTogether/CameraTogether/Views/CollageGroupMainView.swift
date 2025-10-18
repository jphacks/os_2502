import FirebaseAuth
import SwiftUI

struct CollageGroupMainView: View {
    let authManager: AuthenticationManager
    @Environment(\.appColors) var appColors
    @State private var isAnimating = false
    @State private var collageGroupViewModel: CollageGroupViewModel?
    @State private var showCreateGroup = false
    @State private var showJoinGroup = false
    @State private var showWaitingRoom = false
    @State private var selectedGroupType: GroupType = .temporaryLocal
    @State private var groupName: String = ""

    var body: some View {
        ZStack {
            appColors.backgroundGradient
                .ignoresSafeArea()

            VStack(spacing: 0) {
                // ヘッダーセクション
                headerSection
                    .padding(.top, 20)
                    .padding(.horizontal, 24)

                // 空の状態セクション（常に表示）
                emptyStateSection

                Spacer()

                // ボタンセクション
                actionButtonsSection
                    .padding(.horizontal, 24)
                    .padding(.bottom, 32)
            }
        }
        .onAppear {
            // CollageGroupViewModelを初期化
            if collageGroupViewModel == nil {
                collageGroupViewModel = CollageGroupViewModel(authManager: authManager)
            }

            withAnimation(.easeOut(duration: 0.6)) {
                isAnimating = true
            }
        }
        .sheet(isPresented: $showCreateGroup) {
            createGroupSheet
        }
        .sheet(isPresented: $showJoinGroup) {
            NavigationStack {
                JoinGroupView(authManager: authManager)
            }
        }
        .fullScreenCover(isPresented: $showWaitingRoom) {
            if let viewModel = collageGroupViewModel {
                NavigationStack {
                    SimpleWaitingRoomView(viewModel: viewModel)
                }
            }
        }
    }

    // ヘッダーセクション
    private var headerSection: some View {
        VStack(spacing: 12) {
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text("グループコラージュ")
                        .font(.system(size: 28, weight: .bold, design: .rounded))
                        .foregroundColor(appColors.textPrimary)

                    Text("友達と一緒に写真を撮ろう")
                        .font(.system(size: 14, weight: .medium))
                        .foregroundColor(appColors.textSecondary)
                }
                Spacer()
            }
        }
    }

    // 空の状態セクション
    private var emptyStateSection: some View {
        VStack(spacing: 24) {
            Spacer()

            ZStack {
                Circle()
                    .fill(
                        LinearGradient(
                            gradient: Gradient(colors: [
                                Color.blue.opacity(0.2),
                                Color.purple.opacity(0.1),
                            ]),
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(width: 120, height: 120)
                    .scaleEffect(isAnimating ? 1.0 : 0.9)

                Image(systemName: "photo.stack.fill")
                    .font(.system(size: 50))
                    .foregroundStyle(
                        LinearGradient(
                            gradient: Gradient(colors: [Color.blue, Color.purple]),
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
            }
            .shadow(color: Color.blue.opacity(0.3), radius: 20, x: 0, y: 10)

            Text("グループがありません")
                .font(.system(size: 20, weight: .semibold))
                .foregroundColor(appColors.textPrimary)

            Text("新しいグループを作成するか\n既存のグループに参加しましょう")
                .font(.system(size: 14, weight: .medium))
                .foregroundColor(appColors.textSecondary)
                .multilineTextAlignment(.center)
                .lineSpacing(4)

            Spacer()
        }
    }

    // アクションボタンセクション
    private var actionButtonsSection: some View {
        VStack(spacing: 12) {
            Button {
                showCreateGroup = true
            } label: {
                HStack(spacing: 12) {
                    Image(systemName: "plus.circle.fill")
                        .font(.title2)
                    Text("グループを作成")
                        .font(.system(size: 18, weight: .semibold))
                }
                .foregroundColor(.white)
                .frame(maxWidth: .infinity)
                .frame(height: 56)
                .background(
                    LinearGradient(
                        gradient: Gradient(colors: [Color.blue, Color.blue.opacity(0.8)]),
                        startPoint: .leading,
                        endPoint: .trailing
                    )
                )
                .cornerRadius(16)
                .shadow(color: Color.blue.opacity(0.4), radius: 15, x: 0, y: 8)
            }
            .buttonStyle(ScaleButtonStyle())

            Button {
                showJoinGroup = true
            } label: {
                HStack(spacing: 12) {
                    Image(systemName: "person.badge.plus.fill")
                        .font(.title2)
                    Text("グループに参加")
                        .font(.system(size: 18, weight: .semibold))
                }
                .foregroundColor(.blue)
                .frame(maxWidth: .infinity)
                .frame(height: 56)
                .background(
                    RoundedRectangle(cornerRadius: 16)
                        .fill(Color.white)
                        .shadow(color: Color.black.opacity(0.05), radius: 10, x: 0, y: 4)
                )
                .overlay(
                    RoundedRectangle(cornerRadius: 16)
                        .stroke(Color.blue.opacity(0.3), lineWidth: 1.5)
                )
            }
            .buttonStyle(ScaleButtonStyle())
        }
    }

    // グループ作成シート
    private var createGroupSheet: some View {
        NavigationStack {
            ZStack {
                appColors.backgroundGradient
                    .ignoresSafeArea()

                VStack(spacing: 24) {
                    // グループ名入力
                    VStack(alignment: .leading, spacing: 8) {
                        Text("グループ名")
                            .font(.headline)
                            .foregroundColor(appColors.textPrimary)
                        TextField("例: 友達グループ", text: $groupName)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                            .padding(.horizontal, 4)
                    }
                    .padding(.horizontal, 16)

                    VStack(spacing: 16) {
                        GroupTypeButton(
                            title: "ローカル",
                            description: "近くにいる友達と簡単にグループを作成",
                            icon: "location.fill",
                            isSelected: selectedGroupType == .temporaryLocal
                        ) {
                            selectedGroupType = .temporaryLocal
                        }

                        GroupTypeButton(
                            title: "グローバル",
                            description: "インターネットを通じて友達とグループを作成",
                            icon: "globe",
                            isSelected: selectedGroupType == .temporaryGlobal
                        ) {
                            selectedGroupType = .temporaryGlobal
                        }

                        GroupTypeButton(
                            title: "固定",
                            description: "いつでも参加できる固定グループを作成",
                            icon: "lock.fill",
                            isSelected: selectedGroupType == .fixed
                        ) {
                            selectedGroupType = .fixed
                        }
                    }
                    .padding(.horizontal, 16)

                    Spacer()

                    // エラーメッセージ
                    if let errorMessage = collageGroupViewModel?.errorMessage {
                        Text(errorMessage)
                            .font(.caption)
                            .foregroundColor(.red)
                            .multilineTextAlignment(.center)
                            .padding(.horizontal, 16)
                    }

                    Button {
                        Task {
                            await createAndNavigateToGroup()
                        }
                    } label: {
                        if collageGroupViewModel?.isLoading == true {
                            ProgressView()
                                .progressViewStyle(CircularProgressViewStyle(tint: .white))
                                .frame(maxWidth: .infinity)
                                .frame(height: 56)
                        } else {
                            Text("グループを作成")
                                .font(.headline)
                                .foregroundColor(.white)
                                .frame(maxWidth: .infinity)
                                .frame(height: 56)
                        }
                    }
                    .background(
                        LinearGradient(
                            colors:
                                groupName.isEmpty
                                ? [Color.gray, Color.gray.opacity(0.8)]
                                : [Color.blue, Color.blue.opacity(0.8)],
                            startPoint: .leading,
                            endPoint: .trailing
                        )
                    )
                    .cornerRadius(16)
                    .disabled(groupName.isEmpty || collageGroupViewModel?.isLoading == true)
                    .padding(.horizontal, 16)
                    .padding(.bottom, 32)
                }
                .padding(.top, 20)
            }
            .navigationTitle("グループタイプを選択")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("閉じる") {
                        showCreateGroup = false
                    }
                }
            }
        }
    }

    private func createAndNavigateToGroup() async {
        let name = groupName.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !name.isEmpty else { return }
        guard let viewModel = collageGroupViewModel else { return }

        // API経由でグループ作成
        await viewModel.createGroup(
            type: selectedGroupType,
            name: name,
            maxMembers: 10
        )

        // 作成成功したらシートを閉じて待機室に遷移
        if viewModel.currentGroup != nil {
            showCreateGroup = false
            groupName = ""  // 入力をクリア

            // 待機室に遷移
            DispatchQueue.main.asyncAfter(deadline: .now() + 0.3) {
                showWaitingRoom = true
            }
        }
    }
}

// カスタムボタンスタイル
struct ScaleButtonStyle: ButtonStyle {
    func makeBody(configuration: ButtonStyleConfiguration) -> some View {
        configuration.label
            .scaleEffect(configuration.isPressed ? 0.96 : 1.0)
            .animation(.easeInOut(duration: 0.1), value: configuration.isPressed)
    }
}

#Preview {
    NavigationStack {
        CollageGroupMainView(authManager: AuthenticationManager())
    }
}
