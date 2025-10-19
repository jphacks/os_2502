import SwiftUI

struct AddMemberSheetView: View {
    @Environment(\.dismiss) private var dismiss
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors
    let groupType: GroupType
    let onShowQR: () -> Void
    let onFriendSelect: () -> Void

    var body: some View {
        ZStack {
            appColors.backgroundGradient
                .ignoresSafeArea()

            VStack(spacing: 20) {
                HStack {
                    Text("メンバーを追加")
                        .font(.headline)
                        .foregroundColor(appColors.textPrimary)
                    Spacer()
                    Button {
                        dismiss()
                    } label: {
                        Image(systemName: "xmark.circle.fill")
                            .font(.title3)
                            .foregroundColor(appColors.textSecondary)
                    }
                }
                .padding(.horizontal, 24)
                .padding(.top, 20)

                VStack(spacing: 12) {
                    // QRコード表示ボタン（ローカル・グローバル共通）
                    Button {
                        onShowQR()
                    } label: {
                        HStack(spacing: 12) {
                            ZStack {
                                Circle()
                                    .fill(
                                        LinearGradient(
                                            colors: [.blue.opacity(0.6), .cyan.opacity(0.4)],
                                            startPoint: .topLeading,
                                            endPoint: .bottomTrailing
                                        )
                                    )
                                    .frame(width: 40, height: 40)

                                Image(systemName: "qrcode")
                                    .font(.system(size: 18))
                                    .foregroundColor(.white)
                            }

                            Text("QRコード表示")
                                .font(.headline)
                                .foregroundColor(appColors.textPrimary)

                            Spacer()
                        }
                        .padding(16)
                        .glassMorphism(
                            cornerRadius: 16, opacity: colorScheme == .dark ? 0.2 : 0.7)
                    }
                    .buttonStyle(.plain)

                    // フレンド選択ボタン
                    Button {
                        onFriendSelect()
                    } label: {
                        HStack(spacing: 12) {
                            ZStack {
                                Circle()
                                    .fill(
                                        LinearGradient(
                                            colors: [.purple.opacity(0.6), .pink.opacity(0.4)],
                                            startPoint: .topLeading,
                                            endPoint: .bottomTrailing
                                        )
                                    )
                                    .frame(width: 40, height: 40)

                                Image(systemName: "person.2.fill")
                                    .font(.system(size: 18))
                                    .foregroundColor(.white)
                            }

                            Text("フレンドから選択")
                                .font(.headline)
                                .foregroundColor(appColors.textPrimary)

                            Spacer()
                        }
                        .padding(16)
                        .glassMorphism(cornerRadius: 16, opacity: colorScheme == .dark ? 0.2 : 0.7)
                    }
                    .buttonStyle(.plain)
                }
                .padding(.horizontal, 24)

                Spacer()
            }
        }
    }
}
