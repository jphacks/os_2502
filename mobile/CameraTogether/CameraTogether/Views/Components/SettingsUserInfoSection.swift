import SwiftUI

struct SettingsUserInfoSection: View {
    @Environment(\.colorScheme) var colorScheme
    let userName: String
    let iconName: String

    var body: some View {
        HStack(spacing: 16) {
            ZStack {
                Circle()
                    .fill(
                        LinearGradient(
                            colors: [.blue.opacity(0.6), .cyan.opacity(0.4)],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(width: 60, height: 60)

                Image(iconName)
                    .resizable()
                    .scaledToFill()
                    .frame(width: 56, height: 56)
                    .clipShape(Circle())
                    .overlay(
                        Circle()
                            .stroke(
                                LinearGradient(
                                    colors: [
                                        Color.white.opacity(0.5),
                                        Color.white.opacity(0.1),
                                    ],
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                ),
                                lineWidth: 2
                            )
                    )
            }

            Text(userName)
                .font(.system(size: 24, weight: .semibold))
                .foregroundColor(textColor)

            Spacer()
        }
        .padding(.horizontal, 24)
        .padding(.vertical, 20)
        .glassMorphism(cornerRadius: 24, opacity: colorScheme == .dark ? 0.2 : 0.7)
    }

    private var textColor: Color {
        colorScheme == .dark ? .white : .primary
    }
}
