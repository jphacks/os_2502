import SwiftUI

struct SettingsSheetView: View {
    @Environment(\.dismiss) private var dismiss
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors
    @AppStorage("isNotificationEnabled") private var isNotificationEnabled = true

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
                    // ログアウト処理
                }
            )

            SettingsActionButton(
                icon: "trash",
                title: "アカウント削除",
                action: {
                    // アカウント削除処理
                }
            )
        }
    }
}

#Preview {
    SettingsSheetView()
}
