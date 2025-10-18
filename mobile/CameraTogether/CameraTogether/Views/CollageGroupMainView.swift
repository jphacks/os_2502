import SwiftUI

struct CollageGroupMainView: View {
    @Environment(\.appColors) var appColors
    @State private var isAnimating = false

    var body: some View {
        ZStack {
            // グラデーション背景
            appColors.backgroundGradient
                .ignoresSafeArea()

            VStack(spacing: 0) {
                Spacer()

                // ヘッダーセクション
                VStack(spacing: 24) {
                    // アイコン
                    ZStack {
                        Circle()
                            .fill(
                                LinearGradient(
                                    gradient: Gradient(colors: [
                                        Color.blue.opacity(0.2),
                                        Color.purple.opacity(0.1),
                                    ]),
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            )
                            .frame(width: 140, height: 140)
                            .scaleEffect(isAnimating ? 1.0 : 0.9)

                        Image(systemName: "photo.stack.fill")
                            .font(.system(size: 60))
                            .foregroundStyle(
                                LinearGradient(
                                    gradient: Gradient(colors: [Color.blue, Color.purple]),
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            )
                    }
                    .shadow(color: Color.blue.opacity(0.3), radius: 20, x: 0, y: 10)

                    VStack(spacing: 12) {
                        Text("グループコラージュ")
                            .font(.system(size: 32, weight: .bold, design: .rounded))
                            .foregroundStyle(
                                LinearGradient(
                                    gradient: Gradient(colors: [Color.primary, Color.blue]),
                                    startPoint: .leading,
                                    endPoint: .trailing
                                )
                            )

                        Text("友達と一緒に写真を撮って\n素敵なコラージュを作成しよう")
                            .font(.system(size: 16, weight: .medium))
                            .foregroundColor(.secondary)
                            .multilineTextAlignment(.center)
                            .lineSpacing(4)
                    }
                }
                .padding(.top, 60)

                Spacer()

                // ボタンセクション
                VStack(spacing: 16) {
                    NavigationLink {
                        SimpleCreateGroupView()
                    } label: {
                        HStack(spacing: 12) {
                            Image(systemName: "plus.circle.fill")
                                .font(.title2)
                            Text("グループを作成")
                                .font(.system(size: 18, weight: .semibold))
                        }
                        .foregroundColor(.white)
                        .frame(maxWidth: .infinity)
                        .frame(height: 56)
                        .background(
                            LinearGradient(
                                gradient: Gradient(colors: [Color.blue, Color.blue.opacity(0.8)]),
                                startPoint: .leading,
                                endPoint: .trailing
                            )
                        )
                        .cornerRadius(16)
                        .shadow(color: Color.blue.opacity(0.4), radius: 15, x: 0, y: 8)
                    }
                    .buttonStyle(ScaleButtonStyle())

                    NavigationLink {
                        JoinGroupView()
                    } label: {
                        HStack(spacing: 12) {
                            Image(systemName: "person.badge.plus.fill")
                                .font(.title2)
                            Text("グループに参加")
                                .font(.system(size: 18, weight: .semibold))
                        }
                        .foregroundColor(.blue)
                        .frame(maxWidth: .infinity)
                        .frame(height: 56)
                        .background(
                            RoundedRectangle(cornerRadius: 16)
                                .fill(Color.white)
                                .shadow(color: Color.black.opacity(0.05), radius: 10, x: 0, y: 4)
                        )
                        .overlay(
                            RoundedRectangle(cornerRadius: 16)
                                .stroke(Color.blue.opacity(0.3), lineWidth: 1.5)
                        )
                    }
                    .buttonStyle(ScaleButtonStyle())
                }
                .padding(.horizontal, 32)
                .padding(.bottom, 60)
            }
        }
        .onAppear {
            withAnimation(.easeOut(duration: 0.6)) {
                isAnimating = true
            }
        }
    }
}

// カスタムボタンスタイル
struct ScaleButtonStyle: ButtonStyle {
    func makeBody(configuration: ButtonStyleConfiguration) -> some View {
        configuration.label
            .scaleEffect(configuration.isPressed ? 0.96 : 1.0)
            .animation(.easeInOut(duration: 0.1), value: configuration.isPressed)
    }
}

#Preview {
    NavigationStack {
        CollageGroupMainView()
    }
}
