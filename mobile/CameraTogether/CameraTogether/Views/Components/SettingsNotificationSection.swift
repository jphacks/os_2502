import SwiftUI

struct SettingsNotificationSection: View {
    @Environment(\.colorScheme) var colorScheme
    @Binding var isEnabled: Bool

    var body: some View {
        HStack(spacing: 16) {
            ZStack {
                Circle()
                    .fill(
                        LinearGradient(
                            colors: [.cyan.opacity(0.6), .blue.opacity(0.4)],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(width: 40, height: 40)

                Image(systemName: "bell.fill")
                    .font(.system(size: 18))
                    .foregroundColor(.white)
            }

            Text("通知")
                .font(.system(size: 18, weight: .medium))
                .foregroundColor(textColor)

            Spacer()

            Toggle("", isOn: $isEnabled)
                .labelsHidden()
                .tint(.cyan)
        }
        .padding(.horizontal, 20)
        .padding(.vertical, 16)
        .glassMorphism(cornerRadius: 20, opacity: colorScheme == .dark ? 0.2 : 0.7)
    }

    private var textColor: Color {
        colorScheme == .dark ? .white : .primary
    }
}
