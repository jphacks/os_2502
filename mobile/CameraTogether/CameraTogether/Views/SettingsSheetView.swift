import SwiftUI

struct SettingsSheetView: View {
    @Environment(\.dismiss) private var dismiss
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors
    @AppStorage("isNotificationEnabled") private var isNotificationEnabled = true
    var authManager: AuthenticationManager
    @State private var showLogoutAlert = false
    @State private var showDeleteAccountAlert = false
    @State private var errorMessage: String?
    @State private var isLoading = false

    var body: some View {
        NavigationStack {
            ZStack {
                appColors.backgroundGradient
                    .ignoresSafeArea()

                ScrollView {
                    VStack(spacing: 24) {
                        Spacer()
                            .frame(height: 8)

                        userInfoSection
                        notificationSection

                        Spacer()
                            .frame(height: 240)

                        accountActionsSection

                        Spacer()
                            .frame(height: 40)
                    }
                    .padding(.horizontal, 24)
                }
            }
            .navigationTitle("設定")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button {
                        dismiss()
                    } label: {
                        Image(systemName: "xmark.circle.fill")
                            .font(.title3)
                            .foregroundColor(appColors.textPrimary)
                    }
                }
            }
            .toolbarBackground(.visible, for: .navigationBar)
            .toolbarBackground(Color.clear, for: .navigationBar)
        }
        .presentationDetents([.large])
        .presentationDragIndicator(.visible)
    }

    private var userInfoSection: some View {
        // TODO: 結合
        SettingsUserInfoSection(
            userName: "noonyuu",
            iconName: "collage-icon"
        )
    }

    private var notificationSection: some View {
        SettingsNotificationSection(isEnabled: $isNotificationEnabled)
    }

    private var accountActionsSection: some View {
        VStack(spacing: 16) {
            SettingsActionButton(
                icon: "rectangle.portrait.and.arrow.right",
                title: "ログアウト",
                action: {
                    showLogoutAlert = true
                }
            )
            .disabled(isLoading)

            SettingsActionButton(
                icon: "trash",
                title: "アカウント削除",
                action: {
                    showDeleteAccountAlert = true
                }
            )
            .disabled(isLoading)

            if let errorMessage = errorMessage {
                Text(errorMessage)
                    .font(.caption)
                    .foregroundColor(.red)
                    .multilineTextAlignment(.center)
                    .padding(.top, 8)
            }
        }
        .alert("ログアウト", isPresented: $showLogoutAlert) {
            Button("キャンセル", role: .cancel) {}
            Button("ログアウト", role: .destructive) {
                handleLogout()
            }
        } message: {
            Text("ログアウトしますか？")
        }
        .alert("アカウント削除", isPresented: $showDeleteAccountAlert) {
            Button("キャンセル", role: .cancel) {}
            Button("削除", role: .destructive) {
                handleDeleteAccount()
            }
        } message: {
            Text("アカウントを削除すると、すべてのデータが失われます。この操作は取り消せません。")
        }
    }

    private func handleLogout() {
        isLoading = true
        errorMessage = nil

        do {
            try authManager.signOut()
            dismiss()
        } catch {
            errorMessage = "ログアウトに失敗しました: \(error.localizedDescription)"
        }

        isLoading = false
    }

    private func handleDeleteAccount() {
        isLoading = true
        errorMessage = nil

        Task {
            do {
                try await authManager.deleteAccount()
                dismiss()
            } catch {
                errorMessage = "アカウント削除に失敗しました: \(error.localizedDescription)"
            }
            isLoading = false
        }
    }
}

#Preview {
    SettingsSheetView(authManager: AuthenticationManager())
}
