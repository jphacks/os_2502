import CoreImage.CIFilterBuiltins
import SwiftUI

struct GroupQRCodeView: View {
    @Environment(\.dismiss) private var dismiss
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors
    let group: CollageGroup

    var body: some View {
        NavigationStack {
            ZStack {
                appColors.backgroundGradient
                    .ignoresSafeArea()

                VStack(spacing: 24) {
                    Spacer()

                    VStack(spacing: 20) {
                        Text("QRコードをスキャン")
                            .font(.title2)
                            .fontWeight(.bold)
                            .foregroundColor(appColors.textPrimary)

                        Text("このQRコードを他のメンバーに\nスキャンしてもらってください")
                            .font(.subheadline)
                            .foregroundColor(appColors.textSecondary)
                            .multilineTextAlignment(.center)

                        if let qrImage = generateQRCode(from: group.inviteCode) {
                            ZStack {
                                RoundedRectangle(cornerRadius: 24)
                                    .fill(Color.white)
                                    .shadow(
                                        color: Color.black.opacity(0.1), radius: 20, x: 0, y: 10)

                                Image(uiImage: qrImage)
                                    .interpolation(.none)
                                    .resizable()
                                    .scaledToFit()
                                    .frame(width: 250, height: 250)
                            }
                            .frame(width: 280, height: 280)
                        }

                        VStack(spacing: 8) {
                            Text("招待コード")
                                .font(.caption)
                                .foregroundColor(appColors.textSecondary)

                            HStack(spacing: 12) {
                                Text(group.inviteCode)
                                    .font(.system(.title3, design: .monospaced))
                                    .fontWeight(.semibold)
                                    .foregroundColor(appColors.textPrimary)
                                    .padding(.horizontal, 20)
                                    .padding(.vertical, 12)
                                    .glassMorphism(
                                        cornerRadius: 12, opacity: colorScheme == .dark ? 0.2 : 0.7
                                    )

                                Button {
                                    UIPasteboard.general.string = group.inviteCode
                                } label: {
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
                                            .frame(width: 44, height: 44)

                                        Image(systemName: "doc.on.doc.fill")
                                            .font(.system(size: 16))
                                            .foregroundColor(.white)
                                    }
                                }
                            }
                        }
                    }
                    .padding(.horizontal, 24)

                    Spacer()
                }
            }
            .navigationTitle("グループQRコード")
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

    private func generateQRCode(from string: String) -> UIImage? {
        let context = CIContext()
        let filter = CIFilter.qrCodeGenerator()

        filter.message = Data(string.utf8)
        filter.correctionLevel = "M"

        if let outputImage = filter.outputImage {
            let transform = CGAffineTransform(scaleX: 10, y: 10)
            let scaledImage = outputImage.transformed(by: transform)

            if let cgImage = context.createCGImage(scaledImage, from: scaledImage.extent) {
                return UIImage(cgImage: cgImage)
            }
        }

        return nil
    }
}
