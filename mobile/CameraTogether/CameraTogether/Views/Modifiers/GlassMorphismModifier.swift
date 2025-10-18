import SwiftUI

struct GlassMorphismModifier: ViewModifier {
    @Environment(\.colorScheme) var colorScheme

    var cornerRadius: CGFloat = 20
    var opacity: Double = 0.7

    func body(content: Content) -> some View {
        content
            .background(
                ZStack {
                    if colorScheme == .dark {
                        Color.white.opacity(0.1)
                    } else {
                        Color.white.opacity(opacity)
                    }

                    BlurView(style: colorScheme == .dark ? .dark : .light)
                }
            )
            .cornerRadius(cornerRadius)
            .overlay(
                RoundedRectangle(cornerRadius: cornerRadius)
                    .stroke(
                        LinearGradient(
                            colors: [
                                Color.white.opacity(colorScheme == .dark ? 0.3 : 0.5),
                                Color.white.opacity(0.1),
                            ],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        ),
                        lineWidth: 1
                    )
            )
            .shadow(
                color: Color.black.opacity(colorScheme == .dark ? 0.3 : 0.1),
                radius: 10,
                x: 0,
                y: 5
            )
    }
}

struct BlurView: UIViewRepresentable {
    var style: UIBlurEffect.Style

    func makeUIView(context: Context) -> UIVisualEffectView {
        let view = UIVisualEffectView(effect: UIBlurEffect(style: style))
        return view
    }

    func updateUIView(_ uiView: UIVisualEffectView, context: Context) {
        uiView.effect = UIBlurEffect(style: style)
    }
}

extension View {
    func glassMorphism(cornerRadius: CGFloat = 20, opacity: Double = 0.7) -> some View {
        modifier(GlassMorphismModifier(cornerRadius: cornerRadius, opacity: opacity))
    }
}
