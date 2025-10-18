import SwiftUI

// 背景
extension View {
    func appBackgroundGradient(_ colorScheme: ColorScheme) -> some View {
        self.background(
            Group {
                if colorScheme == .dark {
                    LinearGradient(
                        gradient: Gradient(colors: [
                            Color(red: 0.1, green: 0.1, blue: 0.2),
                            Color(red: 0.15, green: 0.1, blue: 0.25),
                            Color(red: 0.2, green: 0.1, blue: 0.3),
                        ]),
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    )
                } else {
                    LinearGradient(
                        gradient: Gradient(colors: [
                            Color.blue.opacity(0.3),
                            Color.purple.opacity(0.3),
                            Color.pink.opacity(0.2),
                        ]),
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    )
                }
            }
        )
    }
}

struct AppColors {
    let colorScheme: ColorScheme

    var textPrimary: Color {
        colorScheme == .dark ? .white : .primary
    }

    var textSecondary: Color {
        colorScheme == .dark ? .white.opacity(0.7) : .secondary
    }

    var textColor: Color {
        colorScheme == .dark ? .white : .primary
    }

    var backgroundOpacity: Double {
        colorScheme == .dark ? 0.2 : 0.7
    }

    var backgroundGradient: some View {
        Group {
            if colorScheme == .dark {
                LinearGradient(
                    gradient: Gradient(colors: [
                        Color(red: 0.1, green: 0.1, blue: 0.2),
                        Color(red: 0.15, green: 0.1, blue: 0.25),
                        Color(red: 0.2, green: 0.1, blue: 0.3),
                    ]),
                    startPoint: .topLeading,
                    endPoint: .bottomTrailing
                )
            } else {
                LinearGradient(
                    gradient: Gradient(colors: [
                        Color.blue.opacity(0.3),
                        Color.purple.opacity(0.3),
                        Color.pink.opacity(0.2),
                    ]),
                    startPoint: .topLeading,
                    endPoint: .bottomTrailing
                )
            }
        }
    }
}

struct AppColorsKey: EnvironmentKey {
    static let defaultValue = AppColors(colorScheme: .light)
}

extension EnvironmentValues {
    var appColors: AppColors {
        get { self[AppColorsKey.self] }
        set { self[AppColorsKey.self] = newValue }
    }
}

extension View {
    func setupAppColors(colorScheme: ColorScheme) -> some View {
        self.environment(\.appColors, AppColors(colorScheme: colorScheme))
    }
}
