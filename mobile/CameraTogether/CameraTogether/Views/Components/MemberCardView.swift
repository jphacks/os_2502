import SwiftUI

struct MemberCardView: View {
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors
    let member: CollageGroupMember

    var body: some View {
        HStack(spacing: 16) {
            ZStack {
                Circle()
                    .fill(
                        LinearGradient(
                            colors: statusGradientColors,
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(width: 50, height: 50)

                Image(systemName: "person.fill")
                    .font(.title3)
                    .foregroundColor(.white)
            }

            VStack(alignment: .leading, spacing: 4) {
                Text(member.name)
                    .font(.headline)
                    .foregroundColor(appColors.textPrimary)

                HStack(spacing: 4) {
                    Image(systemName: statusIcon)
                        .font(.caption)
                    Text(statusText)
                        .font(.caption)
                }
                .foregroundColor(statusColor)
            }

            Spacer()

            if member.isReady {
                ZStack {
                    Circle()
                        .fill(Color.green.opacity(0.2))
                        .frame(width: 32, height: 32)

                    Image(systemName: "checkmark")
                        .font(.system(size: 14, weight: .bold))
                        .foregroundColor(.green)
                }
            }
        }
        .padding(16)
        .glassMorphism(cornerRadius: 16, opacity: colorScheme == .dark ? 0.2 : 0.7)
    }

    private var statusGradientColors: [Color] {
        if member.isReady {
            return [.green.opacity(0.6), .mint.opacity(0.4)]
        } else {
            return [.orange.opacity(0.6), .yellow.opacity(0.4)]
        }
    }

    private var statusIcon: String {
        member.isReady ? "checkmark.circle.fill" : "clock.fill"
    }

    private var statusText: String {
        member.isReady ? "準備完了" : "待機中"
    }

    var statusColor: Color {
        member.isReady ? .green : .orange
    }
}
