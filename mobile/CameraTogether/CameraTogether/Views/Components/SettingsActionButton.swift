import SwiftUI

struct SettingsActionButton: View {
    @Environment(\.colorScheme) var colorScheme
    let icon: String
    let title: String
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            HStack(spacing: 16) {
                ZStack {
                    Circle()
                        .fill(
                            LinearGradient(
                                colors: [.red.opacity(0.7), .pink.opacity(0.5)],
                                startPoint: .topLeading,
                                endPoint: .bottomTrailing
                            )
                        )
                        .frame(width: 40, height: 40)

                    Image(systemName: icon)
                        .font(.system(size: 18))
                        .foregroundColor(.white)
                }

                Text(title)
                    .font(.system(size: 18, weight: .medium))
                    .foregroundColor(.red)

                Spacer()

                Image(systemName: "chevron.right")
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundColor(.red.opacity(0.6))
            }
            .padding(.horizontal, 20)
            .padding(.vertical, 16)
            .glassMorphism(cornerRadius: 20, opacity: colorScheme == .dark ? 0.2 : 0.7)
        }
        .buttonStyle(.plain)
    }
}
