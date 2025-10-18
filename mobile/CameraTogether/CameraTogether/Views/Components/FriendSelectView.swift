import SwiftUI

struct FriendSelectView: View {
    @Environment(\.dismiss) private var dismiss
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors
    let onSelect: (String) -> Void

    private let mockFriends = [
        ("Alice", "person.fill"),
        ("Bob", "person.fill"),
        ("Charlie", "person.fill"),
        ("Diana", "person.fill"),
        ("Eve", "person.fill"),
    ]

    var body: some View {
        NavigationStack {
            ZStack {
                appColors.backgroundGradient
                    .ignoresSafeArea()

                ScrollView {
                    VStack(spacing: 12) {
                        Spacer()
                            .frame(height: 8)

                        ForEach(mockFriends, id: \.0) { friend in
                            Button {
                                onSelect(friend.0)
                            } label: {
                                HStack(spacing: 16) {
                                    ZStack {
                                        Circle()
                                            .fill(
                                                LinearGradient(
                                                    colors: [
                                                        .blue.opacity(0.6), .cyan.opacity(0.4),
                                                    ],
                                                    startPoint: .topLeading,
                                                    endPoint: .bottomTrailing
                                                )
                                            )
                                            .frame(width: 50, height: 50)

                                        Image(systemName: friend.1)
                                            .font(.title3)
                                            .foregroundColor(.white)
                                    }

                                    Text(friend.0)
                                        .font(.headline)
                                        .foregroundColor(appColors.textPrimary)

                                    Spacer()

                                    Image(systemName: "chevron.right")
                                        .font(.system(size: 14, weight: .semibold))
                                        .foregroundColor(appColors.textSecondary)
                                }
                                .padding(16)
                                .glassMorphism(
                                    cornerRadius: 16, opacity: colorScheme == .dark ? 0.2 : 0.7)
                            }
                            .buttonStyle(.plain)
                        }

                        Spacer()
                            .frame(height: 40)
                    }
                    .padding(.horizontal, 24)
                }
            }
            .navigationTitle("フレンドを選択")
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
    }
}
