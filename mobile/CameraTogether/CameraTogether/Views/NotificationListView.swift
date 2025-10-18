import SwiftUI

struct NotificationListView: View {
    @Environment(\.dismiss) private var dismiss
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors

    var body: some View {
        NavigationStack {
            ZStack {
                appColors.backgroundGradient
                    .ignoresSafeArea()

                ScrollView {
                    VStack(spacing: 16) {
                        ForEach(0..<5) { index in
                            NotificationCard(
                                title: "新しい招待",
                                message: "太郎さんがあなたをグループに招待しました",
                                time: "\(index + 1)時間前"
                            )
                        }
                    }
                    .padding()
                }
            }
            .navigationTitle("通知")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("閉じる") {
                        dismiss()
                    }
                    .foregroundColor(appColors.textPrimary)
                }
            }
        }
    }
}

struct NotificationCard: View {
    let title: String
    let message: String
    let time: String
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors

    var body: some View {
        HStack(spacing: 12) {
            ZStack {
                Circle()
                    .fill(
                        LinearGradient(
                            colors: [Color.blue.opacity(0.6), Color.cyan.opacity(0.4)],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(width: 50, height: 50)

                Image(systemName: "bell.fill")
                    .font(.title3)
                    .foregroundColor(.white)
            }

            VStack(alignment: .leading, spacing: 4) {
                Text(title)
                    .font(.headline)
                    .foregroundColor(appColors.textColor)

                Text(message)
                    .font(.subheadline)
                    .foregroundColor(appColors.textColor.opacity(0.7))

                Text(time)
                    .font(.caption)
                    .foregroundColor(appColors.textColor.opacity(0.5))
            }

            Spacer()
        }
        .padding()
        .glassMorphism(cornerRadius: 16, opacity: appColors.backgroundOpacity)
    }
}

#Preview {
    NotificationListView()
}
