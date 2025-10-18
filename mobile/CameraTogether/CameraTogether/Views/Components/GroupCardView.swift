import SwiftUI

struct GroupCardView: View {
    let group: CollageGroup
    @Environment(\.colorScheme) var colorScheme

    var body: some View {
        ZStack {
            RoundedRectangle(cornerRadius: 20)
                .fill(gradientColors)

            VStack(alignment: .leading, spacing: 16) {
                HStack(alignment: .top) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack(spacing: 6) {
                            groupIcon
                            Text(groupTypeText)
                                .font(.system(size: 16, weight: .bold))
                                .foregroundColor(cardTextColor)
                        }

                        Text("\(group.members.count)人参加中")
                            .font(.caption)
                            .foregroundColor(cardTextColor.opacity(0.8))
                    }

                    Spacer()

                    statusBadge
                }

                Spacer()

                HStack(spacing: -12) {
                    ForEach(group.members.prefix(5)) { member in
                        ZStack {
                            Circle()
                                .fill(memberCircleBackground)
                                .frame(width: 40, height: 40)

                            Text(String(member.name.prefix(1)))
                                .font(.system(size: 16, weight: .semibold))
                                .foregroundColor(cardTextColor)
                        }
                        .overlay(
                            Circle()
                                .stroke(memberCircleBorder, lineWidth: 2)
                        )
                    }

                    if group.members.count > 5 {
                        ZStack {
                            Circle()
                                .fill(memberCircleBackground.opacity(0.5))
                                .frame(width: 40, height: 40)

                            Text("+\(group.members.count - 5)")
                                .font(.caption)
                                .fontWeight(.bold)
                                .foregroundColor(cardTextColor)
                        }
                        .overlay(
                            Circle()
                                .stroke(memberCircleBorder, lineWidth: 2)
                        )
                    }
                }
            }
            .padding(20)
        }
        .glassMorphism(cornerRadius: 20, opacity: colorScheme == .dark ? 0.2 : 0.7)
        .frame(height: 160)
    }

    private var gradientColors: LinearGradient {
        let colors: [Color]
        if colorScheme == .dark {
            switch group.status {
            case .recruiting:
                colors = [Color.blue.opacity(0.4), Color.cyan.opacity(0.3)]
            case .readyCheck:
                colors = [Color.orange.opacity(0.4), Color.yellow.opacity(0.3)]
            case .countdown:
                colors = [Color.green.opacity(0.4), Color.mint.opacity(0.3)]
            case .photoTaking:
                colors = [Color.purple.opacity(0.4), Color.pink.opacity(0.3)]
            case .completed:
                colors = [Color.gray.opacity(0.3), Color.gray.opacity(0.2)]
            case .expired:
                colors = [Color.red.opacity(0.3), Color.gray.opacity(0.2)]
            }
        } else {
            switch group.status {
            case .recruiting:
                colors = [Color.blue.opacity(0.6), Color.cyan.opacity(0.4)]
            case .readyCheck:
                colors = [Color.orange.opacity(0.6), Color.yellow.opacity(0.4)]
            case .countdown:
                colors = [Color.green.opacity(0.6), Color.mint.opacity(0.4)]
            case .photoTaking:
                colors = [Color.purple.opacity(0.6), Color.pink.opacity(0.4)]
            case .completed:
                colors = [Color.gray.opacity(0.4), Color.gray.opacity(0.3)]
            case .expired:
                colors = [Color.red.opacity(0.4), Color.gray.opacity(0.3)]
            }
        }
        return LinearGradient(
            colors: colors,
            startPoint: .topLeading,
            endPoint: .bottomTrailing
        )
    }

    private var cardTextColor: Color {
        colorScheme == .dark ? .white : .primary
    }

    private var memberCircleBackground: Color {
        colorScheme == .dark
            ? Color.white.opacity(0.15)
            : Color.white.opacity(0.3)
    }

    private var memberCircleBorder: Color {
        colorScheme == .dark
            ? Color.white.opacity(0.2)
            : Color.white.opacity(0.4)
    }

    private var groupIcon: some View {
        Group {
            switch group.type {
            case .temporaryLocal:
                Image(systemName: "location.fill")
                    .font(.caption)
            case .temporaryGlobal:
                Image(systemName: "globe")
                    .font(.caption)
            case .fixed:
                Image(systemName: "star.fill")
                    .font(.caption)
            }
        }
        .foregroundColor(cardTextColor)
    }

    private var groupTypeText: String {
        switch group.type {
        case .temporaryLocal:
            return "ローカルグループ"
        case .temporaryGlobal:
            return "グローバルグループ"
        case .fixed:
            return "固定グループ"
        }
    }

    private var statusBadge: some View {
        Group {
            switch group.status {
            case .recruiting:
                HStack(spacing: 4) {
                    Image(systemName: "person.badge.plus")
                        .font(.caption2)
                    Text("募集中")
                        .font(.caption2)
                        .fontWeight(.semibold)
                }
                .foregroundColor(cardTextColor)
                .padding(.horizontal, 10)
                .padding(.vertical, 6)
                .background(badgeBackground)
                .clipShape(Capsule())
            case .readyCheck:
                HStack(spacing: 4) {
                    Image(systemName: "clock.fill")
                        .font(.caption2)
                    Text("準備中")
                        .font(.caption2)
                        .fontWeight(.semibold)
                }
                .foregroundColor(cardTextColor)
                .padding(.horizontal, 10)
                .padding(.vertical, 6)
                .background(badgeBackground)
                .clipShape(Capsule())
            case .countdown:
                HStack(spacing: 4) {
                    Image(systemName: "timer")
                        .font(.caption2)
                    Text("カウントダウン")
                        .font(.caption2)
                        .fontWeight(.semibold)
                }
                .foregroundColor(cardTextColor)
                .padding(.horizontal, 10)
                .padding(.vertical, 6)
                .background(badgeBackground)
                .clipShape(Capsule())
            case .photoTaking:
                HStack(spacing: 4) {
                    Image(systemName: "camera.fill")
                        .font(.caption2)
                    Text("撮影中")
                        .font(.caption2)
                        .fontWeight(.semibold)
                }
                .foregroundColor(cardTextColor)
                .padding(.horizontal, 10)
                .padding(.vertical, 6)
                .background(badgeBackground)
                .clipShape(Capsule())
            case .completed:
                HStack(spacing: 4) {
                    Image(systemName: "checkmark.circle.fill")
                        .font(.caption2)
                    Text("完了")
                        .font(.caption2)
                        .fontWeight(.semibold)
                }
                .foregroundColor(cardTextColor.opacity(0.8))
                .padding(.horizontal, 10)
                .padding(.vertical, 6)
                .background(badgeBackground.opacity(0.7))
                .clipShape(Capsule())
            case .expired:
                HStack(spacing: 4) {
                    Image(systemName: "exclamationmark.triangle.fill")
                        .font(.caption2)
                    Text("期限切れ")
                        .font(.caption2)
                        .fontWeight(.semibold)
                }
                .foregroundColor(cardTextColor.opacity(0.8))
                .padding(.horizontal, 10)
                .padding(.vertical, 6)
                .background(badgeBackground.opacity(0.7))
                .clipShape(Capsule())
            }
        }
    }

    private var badgeBackground: Color {
        colorScheme == .dark
            ? Color.white.opacity(0.15)
            : Color.white.opacity(0.3)
    }
}

#Preview {
    var group = CollageGroup(
        type: .temporaryLocal,
        maxMembers: 5,
        ownerId: UUID().uuidString
    )
    group.members = [
        CollageGroupMember(name: "あなた"),
        CollageGroupMember(name: "太郎"),
        CollageGroupMember(name: "花子"),
    ]

    return GroupCardView(group: group)
        .padding()
        .background(Color.gray.opacity(0.1))
}
