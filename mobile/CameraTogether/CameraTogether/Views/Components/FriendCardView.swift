import SwiftUI

struct FriendCardView: View {
    let friend: Friend
    let showActions: Bool
    let onAccept: (() -> Void)?
    let onReject: (() -> Void)?

    init(
        friend: Friend, showActions: Bool = false, onAccept: (() -> Void)? = nil,
        onReject: (() -> Void)? = nil
    ) {
        self.friend = friend
        self.showActions = showActions
        self.onAccept = onAccept
        self.onReject = onReject
    }

    var body: some View {
        HStack(spacing: 16) {
            Image(systemName: friend.iconName)
                .font(.system(size: 32))
                .foregroundColor(.cyan)
                .frame(width: 56, height: 56)
                .background(
                    Circle()
                        .fill(Color.cyan.opacity(0.2))
                )

            VStack(alignment: .leading, spacing: 4) {
                Text(friend.name)
                    .font(.headline)
                    .foregroundColor(.primary)

                if showActions {
                    Text("承認待ち")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
            }

            Spacer()

            if showActions {
                HStack(spacing: 12) {
                    Button {
                        onAccept?()
                    } label: {
                        Image(systemName: "checkmark")
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundColor(.white)
                            .frame(width: 40, height: 40)
                            .background(Color.green)
                            .clipShape(Circle())
                    }

                    Button {
                        onReject?()
                    } label: {
                        Image(systemName: "xmark")
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundColor(.white)
                            .frame(width: 40, height: 40)
                            .background(Color.red)
                            .clipShape(Circle())
                    }
                }
            } else {
                Image(systemName: "chevron.right")
                    .font(.system(size: 14))
                    .foregroundColor(.secondary)
            }
        }
        .padding(16)
        .background(
            RoundedRectangle(cornerRadius: 16)
                .fill(Color("button"))
        )
    }
}
